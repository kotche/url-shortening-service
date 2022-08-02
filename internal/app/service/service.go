package service

import (
	"context"
	"log"
	"time"

	"github.com/kotche/url-shortening-service/internal/app/usecase"
	"github.com/kotche/url-shortening-service/internal/config"
)

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Storage interface {
	Add(userID string, url *usecase.URL) error
	GetByID(id string) (*usecase.URL, error)
	GetUserURLs(userID string) ([]*usecase.URL, error)
	Close() error
}

type Database interface {
	Storage
	Ping() error
	WriteBatch(ctx context.Context, userID string, urls map[string]*usecase.URL) error
	DeleteBatch(ctx context.Context, toDelete []usecase.DeleteUserURLs) error
}

type IGenerator interface {
	MakeShortURL() string
}

type Service struct {
	st           Storage
	db           Database
	gen          IGenerator
	deletionChan chan usecase.DeleteUserURLs
	buf          []usecase.DeleteUserURLs
	timer        *time.Timer
	isTimeout    bool
}

func NewService(st Storage, gen IGenerator) *Service {
	s := Service{
		st:           st,
		gen:          gen,
		deletionChan: make(chan usecase.DeleteUserURLs),
		buf:          make([]usecase.DeleteUserURLs, 0, config.BufLen),
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

func (s *Service) GetURLModel(userID string, originURL string) (*usecase.URL, error) {

	var urlModel *usecase.URL

	for {
		shortURL := s.gen.MakeShortURL()
		urlModel, _ = s.st.GetByID(shortURL)

		if urlModel == nil {
			urlModel = usecase.NewURL(originURL, shortURL)
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

func (s *Service) GetURLModelByID(shortURL string) (*usecase.URL, error) {
	urlModel, err := s.st.GetByID(shortURL)

	if err != nil {
		return nil, err
	}
	return urlModel, nil
}

func (s *Service) GetUserURLs(userID string) ([]*usecase.URL, error) {
	userURLs, err := s.st.GetUserURLs(userID)

	if err != nil {
		return nil, err
	}
	return userURLs, nil
}

func (s *Service) Ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}
	return nil
}

func (s *Service) ShortenBatch(ctx context.Context, userID string, input []usecase.InputCorrelationURL) ([]usecase.OutputCorrelationURL, error) {

	output := make([]usecase.OutputCorrelationURL, 0)
	urls := make(map[string]*usecase.URL)
	for _, correlationURL := range input {
		var urlModel *usecase.URL

		for {
			shortURL := s.gen.MakeShortURL()
			if _, ok := urls[shortURL]; ok {
				continue
			}
			urlModel = usecase.NewURL(correlationURL.Origin, shortURL)
			urls[shortURL] = urlModel
			break
		}

		out := usecase.OutputCorrelationURL{
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
		delUserURLs := usecase.DeleteUserURLs{UserID: userID, Short: v}
		s.deletionChan <- delUserURLs
	}
}

func (s *Service) flush(ctx context.Context) {
	del := make([]usecase.DeleteUserURLs, len(s.buf))
	copy(del, s.buf)
	s.buf = make([]usecase.DeleteUserURLs, 0)
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
