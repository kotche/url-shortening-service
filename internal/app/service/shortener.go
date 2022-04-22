package service

import (
	"math/rand"

	"github.com/kotche/url-shortening-service/internal/config"
)

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func MakeShortURL() string {
	b := make([]byte, config.ShortURLLen)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}
