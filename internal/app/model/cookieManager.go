package model

import (
	"context"
	"net/http"

	"github.com/kotche/url-shortening-service/internal/app/config"
	"github.com/kotche/url-shortening-service/internal/app/utils"
)

type CookieManager struct{}

// GetUserID returns the user ID from the cookie
func (c CookieManager) GetUserID(r *http.Request) string {
	return r.Context().Value(config.UserIDCookieName).(string)
}

type CookieManagerMD struct{}

// GetUserID returns the user ID from the metadata cookie
func (c CookieManagerMD) GetUserID(ctx context.Context) string {
	return utils.GetUserIDFromMD(ctx)
}
