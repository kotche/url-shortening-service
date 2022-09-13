package usecase

import "fmt"

// ConflictURLError called if the URL already exists
type ConflictURLError struct {
	ShortenURL string
}

func (e ConflictURLError) Error() string {
	return fmt.Sprintf("url %v already exists", e.ShortenURL)
}

// GoneError called if the URL has already been deleted
type GoneError struct {
	ShortenURL string
}

func (e GoneError) Error() string {
	return fmt.Sprintf("url %v gone", e.ShortenURL)
}
