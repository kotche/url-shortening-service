package storage

import (
	"fmt"

	"github.com/kotche/url-shortening-service/internal/app/usecase"
)

type URLStorage struct {
	urls      map[string]*usecase.URL
	urlsUsers map[string][]*usecase.URL
}

func NewUrls() *URLStorage {
	return &URLStorage{
		urls:      make(map[string]*usecase.URL),
		urlsUsers: make(map[string][]*usecase.URL),
	}
}

func (m *URLStorage) Add(userID string, url *usecase.URL) error {
	m.urls[url.Short] = url
	m.urlsUsers[userID] = append(m.urlsUsers[userID], url)
	return nil
}

func (m *URLStorage) GetByID(id string) (*usecase.URL, error) {

	original, ok := m.urls[id]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return original, nil
}

func (m *URLStorage) GetUserURLs(userID string) ([]*usecase.URL, error) {
	usersURLs := m.urlsUsers[userID]
	return usersURLs, nil
}

func (m *URLStorage) Close() error {
	return nil
}
