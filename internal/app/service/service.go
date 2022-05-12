package service

import (
	"math/rand"
	"time"

	"github.com/kotche/url-shortening-service/internal/config"
)

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Storage interface {
	Add(url *URL) error
	GetByID(id string) (*URL, error)
	Close() error
}

type Service struct {
	st Storage
}

func NewService(st Storage) *Service {
	return &Service{st: st}
}

func (s *Service) MakeShortURL() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, config.ShortURLLen)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}

func (s *Service) GetURLModel(originURL string) (*URL, error) {
	shortURL := s.MakeShortURL()
	urlModel, _ := s.st.GetByID(shortURL)

	for {
		if urlModel == nil {
			urlModel = NewURL(originURL, shortURL)
			err := s.st.Add(urlModel)
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
