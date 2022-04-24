package service

import (
	"math/rand"
	"time"

	"github.com/kotche/url-shortening-service/internal/config"
)

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func MakeShortURL() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, config.ShortURLLen)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}
