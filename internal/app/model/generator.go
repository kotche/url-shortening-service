package model

import (
	"math/rand"

	"github.com/kotche/url-shortening-service/internal/app/config"
)

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Generator generates shortened RL
type Generator struct{}

// MakeShortURL returns the generated shortened URL
func (g Generator) MakeShortURL() string {
	b := make([]byte, config.ShortURLLen)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}
