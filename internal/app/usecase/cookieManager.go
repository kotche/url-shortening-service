package usecase

import (
	"net/http"

	"github.com/kotche/url-shortening-service/internal/app/config"
)

type CookieManager struct{}

func (c CookieManager) GetUserID(r *http.Request) string {
	return r.Context().Value(config.UserIDCookieName).(string)
}
