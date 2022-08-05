package test

import (
	"context"

	"github.com/kotche/url-shortening-service/internal/app/usecase"
)

type FakeRepo struct {
	Short string
}

func (f *FakeRepo) Add(userID string, url *usecase.URL) error {
	url.Short = f.Short
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

func (f *FakeRepo) Ping() error {
	return nil
}

func (f *FakeRepo) WriteBatch(ctx context.Context, userID string, urls map[string]*usecase.URL) error {
	return nil
}

func (f *FakeRepo) DeleteBatch(ctx context.Context, toDelete []usecase.DeleteUserURLs) error {
	return nil
}
