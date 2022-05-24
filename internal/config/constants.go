package config

const (
	ShortURLLen  = 7
	Compression  = "gzip"
	UserIDCookie = "userID"
	UserIDLen    = 16
)

func GetSecretKey() []byte {
	return []byte("be55d1079e6c6167118ac91318fe")
}
