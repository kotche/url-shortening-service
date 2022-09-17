package test

import (
	"context"

	"github.com/kotche/url-shortening-service/internal/app/model"
)

type FakeRepo struct {
	Short string
}

func (f *FakeRepo) Add(userID string, url *model.URL) error {
	url.Short = f.Short
	return nil
}

func (f *FakeRepo) GetByID(id string) (*model.URL, error) {
	return nil, nil
}

func (f *FakeRepo) GetUserURLs(userID string) ([]*model.URL, error) {
	return nil, nil
}

func (f *FakeRepo) Close() error {
	return nil
}

func (f *FakeRepo) Ping(ctx context.Context) error {
	return nil
}

func (f *FakeRepo) WriteBatch(ctx context.Context, userID string, urls map[string]*model.URL) error {
	return nil
}

func (f *FakeRepo) DeleteBatch(ctx context.Context, toDelete []model.DeleteUserURLs) error {
	return nil
}

func (f *FakeRepo) GetNumberOfURLs(ctx context.Context) (int, error) {
	return 0, nil
}

func (f *FakeRepo) GetNumberOfUsers(ctx context.Context) (int, error) {
	return 0, nil
}
