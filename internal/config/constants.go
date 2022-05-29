package config

type ContextType string

const (
	ShortURLLen                  = 7
	Compression                  = "gzip"
	UserIDCookieName ContextType = "user_id"
	CookieMaxAge                 = 86400
	UserIDLen                    = 16
)

func GetSecretKey() []byte {
	return []byte("be55d1079e6c6167118ac91318fe")
}
