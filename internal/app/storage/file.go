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
	file    *os.File
	encoder *json.Encoder
	urls    map[string]*service.URL
}

func NewFileStorage(fileName string) (*FileStorage, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	urls := make(map[string]*service.URL)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		model := &service.URL{}
		err := json.Unmarshal(scanner.Bytes(), model)
		if err != nil {
			return nil, err
		}
		urls[model.Short] = model
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &FileStorage{
		file:    file,
		encoder: json.NewEncoder(file),
		urls:    urls,
	}, nil
}

func (f *FileStorage) Add(url *service.URL) error {
	mu := &sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	f.urls[url.Short] = url

	err := f.encoder.Encode(url)
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

func (f *FileStorage) Close() error {
	return f.file.Close()
}
