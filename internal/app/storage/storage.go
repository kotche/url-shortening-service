package storage

import (
	"fmt"
	"github.com/kotche/url-shortening-service/internal/app/service"
)

type Storage interface {
	Add(url *service.URL)
	GetById(id string) (*service.URL, error)
}

var _ Storage = (*UrlsStorage)(nil)

type UrlsStorage struct {
	urls map[string]*service.URL
}

func NewUrls() *UrlsStorage {
	return &UrlsStorage{
		urls: make(map[string]*service.URL),
	}
}

func (m *UrlsStorage) Add(url *service.URL) {
	m.urls[url.GetShort()] = url
}

func (m *UrlsStorage) GetById(id string) (*service.URL, error) {

	original, ok := m.urls[id]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return original, nil
}
