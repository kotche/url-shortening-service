package service

import (
	"math/rand"
	"time"

	"github.com/kotche/url-shortening-service/internal/config"
)

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Storage interface {
	Add(userID string, url *URL) error
	GetByID(id string) (*URL, error)
	GetUserURLs(userID string) ([]*URL, error)
	Close() error
}

type Database interface {
	Ping() error
}

type Service struct {
	st Storage
	db Database
}

func NewService(st Storage) *Service {
	return &Service{st: st}
}

func (s *Service) SetDB(db Database) {
	s.db = db
}

func (s *Service) MakeShortURL() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, config.ShortURLLen)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}

func (s *Service) GetURLModel(userID string, originURL string) (*URL, error) {
	shortURL := s.MakeShortURL()
	urlModel, _ := s.st.GetByID(shortURL)

	for {
		if urlModel == nil {
			urlModel = NewURL(originURL, shortURL)
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

func (s *Service) GetURLModelByID(shortURL string) (*URL, error) {
	urlModel, err := s.st.GetByID(shortURL)

	if err != nil {
		return nil, err
	}
	return urlModel, nil
}

func (s *Service) GetUserURLs(userID string) ([]*URL, error) {
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
