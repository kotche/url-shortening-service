package usecase

import "fmt"

type ErrConflictURL struct {
	Err        error
	ShortenURL string
}

func (e ErrConflictURL) Error() string {
	return fmt.Sprintf("url %v already exists", e.ShortenURL)
}
