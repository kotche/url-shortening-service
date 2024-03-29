package middlewares

import (
	"context"
	"net/http"

	"github.com/kotche/url-shortening-service/internal/app/config"
	"github.com/kotche/url-shortening-service/internal/app/utils"
)

// UserCookieHandler checks for the presence of the user ID in the cookie file. If not, then a new one is issued
func UserCookieHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			userID       string
			userIDCookie string
		)

		cookieID := utils.GetCookieParam(r, string(config.UserIDCookieName))
		if cookieID != "" {
			userID = utils.GetUserIDFromCookie(cookieID)
			if userID != "" {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), config.UserIDCookieName, userID)))
				return
			}
		}

		userID, userIDCookie = utils.MakeUserIDCookie()
		cookie := http.Cookie{Name: string(config.UserIDCookieName), Value: userIDCookie, Path: "/", MaxAge: config.CookieMaxAge}
		http.SetCookie(w, &cookie)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), config.UserIDCookieName, userID)))
	})
}
