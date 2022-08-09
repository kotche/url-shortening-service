package mock

import (
	"net/http"
)

type CookieManager struct {
	Cookie string
}

func (c CookieManager) GetUserID(r *http.Request) string {
	return c.Cookie
}
