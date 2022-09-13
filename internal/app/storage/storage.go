package storage

import (
	"fmt"

	"github.com/kotche/url-shortening-service/internal/app/model"
)

// URLStorage store the URL in RAM
type URLStorage struct {
	urls      map[string]*model.URL
	urlsUsers map[string][]*model.URL
}

func NewUrls() *URLStorage {
	return &URLStorage{
		urls:      make(map[string]*model.URL),
		urlsUsers: make(map[string][]*model.URL),
	}
}

func (m *URLStorage) Add(userID string, url *model.URL) error {
	m.urls[url.Short] = url
	m.urlsUsers[userID] = append(m.urlsUsers[userID], url)
	return nil
}

func (m *URLStorage) GetByID(id string) (*model.URL, error) {

	original, ok := m.urls[id]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return original, nil
}

func (m *URLStorage) GetUserURLs(userID string) ([]*model.URL, error) {
	usersURLs := m.urlsUsers[userID]
	return usersURLs, nil
}

func (m *URLStorage) Close() error {
	return nil
}
