package storage

import (
	"fmt"

	"github.com/kotche/url-shortening-service/internal/app/service"
)

type urlsStorage struct {
	urls map[string]*service.URL
}

func NewUrls() *urlsStorage {
	return &urlsStorage{
		urls: make(map[string]*service.URL),
	}
}

func (m *urlsStorage) Add(url *service.URL) {
	m.urls[url.GetShort()] = url
}

func (m *urlsStorage) GetByID(id string) (*service.URL, error) {

	original, ok := m.urls[id]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return original, nil
}
