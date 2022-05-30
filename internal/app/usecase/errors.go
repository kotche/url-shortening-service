package usecase

import "fmt"

type ConflictURLError struct {
	Err        error
	ShortenURL string
}

func (e ConflictURLError) Error() string {
	return fmt.Sprintf("url %v already exists", e.ShortenURL)
}
