package config

type ContextType string

const (
	ShortURLLen                  = 7
	Compression                  = "gzip"
	UserIDCookieName ContextType = "userID"
	UserIDLen                    = 16
)

func GetSecretKey() []byte {
	return []byte("be55d1079e6c6167118ac91318fe")
}
