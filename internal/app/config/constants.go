package config

type ContextType string

const (
	ShortURLLen                  = 7
	Compression                  = "gzip"
	UserIDCookieName ContextType = "user_id"
	CookieMaxAge                 = 86400
	UserIDLen                    = 8
	BufLen                       = 3
	Timeout                      = 5
)

// GetSecretKey returns the secret key for generating the encrypted user id in cookies
func GetSecretKey() []byte {
	return []byte("be55d1079e6c6167118ac91318fe")
}
