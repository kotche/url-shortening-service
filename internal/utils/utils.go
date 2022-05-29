package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/kotche/url-shortening-service/internal/config"
)

func GetCookieParam(r *http.Request, name string) string {
	cookieParam, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookieParam.Value
}

func MakeUserIDCookie() (string, string) {
	userID := make([]byte, config.UserIDLen)

	rand.Seed(time.Now().UnixNano())
	rand.Read(userID)
	encodedID := hex.EncodeToString(userID)

	h := hmac.New(sha256.New, config.GetSecretKey())
	h.Write(userID)
	hash := h.Sum(nil)
	return encodedID, encodedID + hex.EncodeToString(hash)
}

func GetUserIDFromCookie(CookieID string) string {
	data, err := hex.DecodeString(CookieID)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	id := data[:config.UserIDLen]
	h := hmac.New(sha256.New, config.GetSecretKey())
	h.Write(id)
	hash := h.Sum(nil)

	if hmac.Equal(hash, data[config.UserIDLen:]) {
		return hex.EncodeToString(id)
	} else {
		log.Printf("UserID %d no auth", id)
		return ""
	}
}
