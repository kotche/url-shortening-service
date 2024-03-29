package utils

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math/rand"
	"net/http"

	"github.com/kotche/url-shortening-service/internal/app/config"
	"google.golang.org/grpc/metadata"
)

// GetCookieParam returns the cookie parameter by name
func GetCookieParam(r *http.Request, name string) string {
	cookieParam, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookieParam.Value
}

// MakeUserIDCookie generates an encrypted user id for cookies
func MakeUserIDCookie() (string, string) {
	userID := make([]byte, config.UserIDLen)

	rand.Read(userID)
	encodedID := hex.EncodeToString(userID)

	h := hmac.New(sha256.New, config.GetSecretKey())
	h.Write(userID)
	hash := h.Sum(nil)
	return encodedID, encodedID + hex.EncodeToString(hash)
}

// GetUserIDFromCookie receives an encrypted user id from the cookie
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
		log.Printf("UserID %v no auth", string(id))
		return ""
	}
}

func GetUserIDFromMD(ctx context.Context) string {
	var (
		userID       string
		userIDCookie string
	)
	userIDName := string(config.UserIDCookieName)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get(userIDName)
		if len(values) > 0 {
			userIDCookie = values[0]
		}
	}
	if userIDCookie != "" {
		userID = GetUserIDFromCookie(userIDCookie)
		if userID != "" {
			return userID
		}
	}
	return ""
}
