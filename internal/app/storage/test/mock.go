package test

import (
	"github.com/kotche/url-shortening-service/internal/app/service"
)

type Mock struct {
	original string
	short    string
	err      bool
}

func NewMock(original, short string) *Mock {
	return &Mock{
		original: original,
		short:    short,
	}
}

func (m *Mock) Add(userID string, url *service.URL) error {
	url.Short = m.short
	return nil
}

func (m *Mock) GetByID(id string) (*service.URL, error) {
	return nil, nil
}

func (m *Mock) GetUserURLs(userID string) ([]*service.URL, error) {
	return nil, nil
}

func (m *Mock) Close() error {
	return nil
}
