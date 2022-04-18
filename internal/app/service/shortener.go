package service

import (
	"github.com/kotche/url-shortening-service/internal/config"
	"math/rand"
)

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func MakeShortURl() string {
	b := make([]byte, config.ShortUrlLen)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}
