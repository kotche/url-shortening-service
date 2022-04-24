package storage

import (
	"fmt"

	"github.com/kotche/url-shortening-service/internal/app/service"
)

type URLStorage struct {
	urls map[string]*service.URL
}

func NewUrls() *URLStorage {
	return &URLStorage{
		urls: make(map[string]*service.URL),
	}
}

func (m *URLStorage) Add(url *service.URL) {
	m.urls[url.GetShort()] = url
}

func (m *URLStorage) GetByID(id string) (*service.URL, error) {

	original, ok := m.urls[id]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return original, nil
}
