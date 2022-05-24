package middlewares

import (
	"context"
	"net/http"

	"github.com/kotche/url-shortening-service/internal/config"
	"github.com/kotche/url-shortening-service/internal/utils"
)

func UserCookieHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			userID       string
			userIDCookie string
		)

		userIDCookieString := config.UserIDCookie

		cookieID := utils.GetCookieParam(r, userIDCookieString)
		if cookieID != "" {
			userID = utils.GetUserIDFromCookie(cookieID)
			if userID != "" {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDCookieString, userID)))
				return
			}
		}

		userID, userIDCookie = utils.MakeUserIDCookie()
		cookie := http.Cookie{Name: config.UserIDCookie, Value: userIDCookie}
		http.SetCookie(w, &cookie)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDCookieString, userID)))
	})
}
