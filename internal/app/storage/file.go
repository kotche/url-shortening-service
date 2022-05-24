package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/kotche/url-shortening-service/internal/app/service"
)

type FileStorage struct {
	file      *os.File
	encoder   *json.Encoder
	urls      map[string]*service.URL
	urlsUsers map[string][]*service.URL
}

type DataFile struct {
	Owner string `json:"owner"`
	*service.URL
}

func NewFileStorage(fileName string) (*FileStorage, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	urls := make(map[string]*service.URL)
	urlsUsers := make(map[string][]*service.URL)

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

func (f *FileStorage) Add(userID string, url *service.URL) error {
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

func (f *FileStorage) GetByID(id string) (*service.URL, error) {
	original, ok := f.urls[id]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	return original, nil
}

func (f *FileStorage) GetUserURLs(userID string) ([]*service.URL, error) {
	usersURLs := f.urlsUsers[userID]
	return usersURLs, nil
}

func (f *FileStorage) Close() error {
	return f.file.Close()
}
