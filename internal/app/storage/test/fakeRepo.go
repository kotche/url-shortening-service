package test

import (
	"github.com/kotche/url-shortening-service/internal/app/usecase"
)

type FakeRepo struct {
	original string
	short    string
	err      bool
}

func NewFakeRepo(original, short string) *FakeRepo {
	return &FakeRepo{
		original: original,
		short:    short,
	}
}

func (f *FakeRepo) Add(userID string, url *usecase.URL) error {
	url.Short = f.short
	return nil
}

func (f *FakeRepo) GetByID(id string) (*usecase.URL, error) {
	return nil, nil
}

func (f *FakeRepo) GetUserURLs(userID string) ([]*usecase.URL, error) {
	return nil, nil
}

func (f *FakeRepo) Close() error {
	return nil
}
