package storage

import (
	"fmt"

	"github.com/kotche/url-shortening-service/internal/app/service"
)

type URLStorage struct {
	urls      map[string]*service.URL
	urlsUsers map[string][]*service.URL
}

func NewUrls() *URLStorage {
	return &URLStorage{
		urls:      make(map[string]*service.URL),
		urlsUsers: make(map[string][]*service.URL),
	}
}

func (m *URLStorage) Add(userID string, url *service.URL) error {
	m.urls[url.Short] = url
	m.urlsUsers[userID] = append(m.urlsUsers[userID], url)
	return nil
}

func (m *URLStorage) GetByID(id string) (*service.URL, error) {

	original, ok := m.urls[id]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return original, nil
}

func (m *URLStorage) GetUserURLs(userID string) ([]*service.URL, error) {
	usersURLs := m.urlsUsers[userID]
	return usersURLs, nil
}

func (m *URLStorage) Close() error {
	return nil
}
