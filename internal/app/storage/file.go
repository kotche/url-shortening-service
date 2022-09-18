package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/kotche/url-shortening-service/internal/app/model"
)

type FileStorage struct {
	file      *os.File
	encoder   *json.Encoder
	urls      map[string]*model.URL
	urlsUsers map[string][]*model.URL
}

// DataFile store the URL in the file system
type DataFile struct {
	Owner string `json:"owner"`
	*model.URL
}

func NewFileStorage(fileName string) (*FileStorage, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	urls := make(map[string]*model.URL)
	urlsUsers := make(map[string][]*model.URL)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dataFile := &DataFile{}
		err = json.Unmarshal(scanner.Bytes(), dataFile)
		if err != nil {
			return nil, err
		}
		urls[dataFile.Short] = dataFile.URL
		urlsUsers[dataFile.Owner] = append(urlsUsers[dataFile.Owner], dataFile.URL)
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return &FileStorage{
		file:      file,
		encoder:   json.NewEncoder(file),
		urls:      urls,
		urlsUsers: urlsUsers,
	}, nil
}

func (f *FileStorage) Add(_ context.Context, userID string, url *model.URL) error {
	mu := &sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	dataFile := &DataFile{}
	dataFile.URL = url
	dataFile.Owner = userID

	f.urls[url.Short] = url
	f.urlsUsers[userID] = append(f.urlsUsers[userID], url)

	err := f.encoder.Encode(dataFile)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileStorage) GetByID(_ context.Context, id string) (*model.URL, error) {
	original, ok := f.urls[id]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return original, nil
}

func (f *FileStorage) GetUserURLs(_ context.Context, userID string) ([]*model.URL, error) {
	usersURLs := f.urlsUsers[userID]
	return usersURLs, nil
}

func (f *FileStorage) Close() error {
	return f.file.Close()
}
