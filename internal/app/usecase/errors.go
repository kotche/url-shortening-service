package usecase

import "fmt"

type ConflictURLError struct {
	ShortenURL string
}

func (e ConflictURLError) Error() string {
	return fmt.Sprintf("url %v already exists", e.ShortenURL)
}

type GoneError struct {
	ShortenURL string
}

func (e GoneError) Error() string {
	return fmt.Sprintf("url %v gone", e.ShortenURL)
}
