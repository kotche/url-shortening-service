package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kotche/url-shortening-service/internal/app/config"
	"github.com/kotche/url-shortening-service/internal/app/model"
	"golang.org/x/sync/errgroup"
)

// Storage describes methods for storage in RAM , file system , and database
type Storage interface {
	Add(userID string, url *model.URL) error
	GetByID(id string) (*model.URL, error)
	GetUserURLs(userID string) ([]*model.URL, error)
	Close() error
}

// Database describes methods for storage in database
type Database interface {
	Storage
	Ping(ctx context.Context) error
	WriteBatch(ctx context.Context, userID string, urls map[string]*model.URL) error
	DeleteBatch(ctx context.Context, toDelete []model.DeleteUserURLs) error
	GetNumberOfUsers(ctx context.Context) (int, error)
	GetNumberOfURLs(ctx context.Context) (int, error)
}

// IGenerator describes methods for generating shortened links
type IGenerator interface {
	MakeShortURL() string
}

type Service struct {
	st           Storage
	db           Database
	Gen          IGenerator
	deletionChan chan model.DeleteUserURLs
	buf          []model.DeleteUserURLs
	timer        *time.Timer
	isTimeout    bool
}

func NewService(st Storage) *Service {
	s := Service{
		st:           st,
		Gen:          model.Generator{},
		deletionChan: make(chan model.DeleteUserURLs),
		buf:          make([]model.DeleteUserURLs, 0, config.BufLen),
		isTimeout:    true,
		timer:        time.NewTimer(0),
	}

	return &s
}

func (s *Service) RunWorker() {
	go s.worker()
}

func (s *Service) SetDB(db Database) {
	s.db = db
}

func (s *Service) GetURLModel(userID string, originURL string) (*model.URL, error) {

	var urlModel *model.URL

	for {
		shortURL := s.Gen.MakeShortURL()
		urlModel, _ = s.st.GetByID(shortURL)

		if urlModel == nil {
			urlModel = model.NewURL(originURL, shortURL)
			err := s.st.Add(userID, urlModel)
			if err != nil {
				return nil, err
			}
		} else if urlModel.Origin != originURL {
			continue
		}
		break
	}

	return urlModel, nil
}

func (s *Service) GetURLModelByID(shortURL string) (*model.URL, error) {
	urlModel, err := s.st.GetByID(shortURL)

	if err != nil {
		return nil, err
	}
	return urlModel, nil
}

func (s *Service) GetUserURLs(userID string) ([]*model.URL, error) {
	userURLs, err := s.st.GetUserURLs(userID)

	if err != nil {
		return nil, err
	}
	return userURLs, nil
}

func (s *Service) Ping(ctx context.Context) error {
	if s.db == nil {
		log.Printf("Ping error: database not initialized")
		return fmt.Errorf("database not initialized")
	}

	if err := s.db.Ping(ctx); err != nil {
		log.Printf("Ping error: %s", err.Error())
		return err
	}
	return nil
}

func (s *Service) ShortenBatch(ctx context.Context, userID string, input []model.InputCorrelationURL) ([]model.OutputCorrelationURL, error) {

	output := make([]model.OutputCorrelationURL, 0, len(input))
	urls := make(map[string]*model.URL)
	for _, correlationURL := range input {
		var urlModel *model.URL

		for {
			shortURL := s.Gen.MakeShortURL()
			if _, ok := urls[shortURL]; ok {
				continue
			}
			urlModel = model.NewURL(correlationURL.Origin, shortURL)
			urls[shortURL] = urlModel
			break
		}

		out := model.OutputCorrelationURL{
			CorrelationID: correlationURL.CorrelationID,
			Short:         urlModel.Short,
		}

		output = append(output, out)
	}

	err := s.db.WriteBatch(ctx, userID, urls)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (s *Service) DeleteURLs(userID string, toDelete []string) {
	for _, v := range toDelete {
		delUserURLs := model.DeleteUserURLs{UserID: userID, Short: v}
		s.deletionChan <- delUserURLs
	}
}

func (s *Service) flush(ctx context.Context) {
	del := make([]model.DeleteUserURLs, len(s.buf))
	copy(del, s.buf)
	s.buf = make([]model.DeleteUserURLs, 0)
	go func() {
		err := s.db.DeleteBatch(ctx, del)
		if err != nil {
			log.Printf("error deleting: " + err.Error())
		}
	}()
}

func (s *Service) worker() {
	ctx := context.Background()

	for {
		select {
		case delRequest := <-s.deletionChan:
			if s.isTimeout {
				s.timer.Reset(time.Second * config.Timeout)
				s.isTimeout = false
			}
			s.buf = append(s.buf, delRequest)
			if len(s.buf) >= config.BufLen {
				s.flush(ctx)
				s.timer.Stop()
				s.isTimeout = true
			}
		case <-s.timer.C:
			if len(s.buf) > 0 {
				s.flush(ctx)
			}
			s.isTimeout = true
		}
	}
}

func (s *Service) GetStats(ctx context.Context) (model.Stats, error) {
	var numberOfURLs, numberOfUsers int
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		res, err := s.db.GetNumberOfURLs(ctx)
		if err != nil {
			log.Printf("%s error: %s", "pgx.GetNumberOfURLs", err.Error())
			return err
		}
		numberOfURLs = res
		return nil
	})

	grp.Go(func() error {
		res, err := s.db.GetNumberOfUsers(ctx)
		if err != nil {
			log.Printf("%s error: %s", "pgx.GetNumberOfUsers", err.Error())
			return err
		}
		numberOfUsers = res
		return nil
	})

	if err := grp.Wait(); err != nil {
		return model.Stats{}, err
	}

	return model.Stats{
		NumberOfURLs:  numberOfURLs,
		NumberOfUsers: numberOfUsers,
	}, nil
}
