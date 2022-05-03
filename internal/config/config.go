package config

import (
	"os"
)

const ShortURLLen = 7

type Config struct {
	serverAddr string
	baseURL    string
}

func NewConfig() *Config {
	return &Config{
		serverAddr: getEnvValue("SERVER_ADDRESS", "localhost:8080"),
		baseURL:    getEnvValue("BASE_URL", "http://localhost:8080"),
	}
}

func (c *Config) GetServerAddr() string {
	return c.serverAddr
}

func (c *Config) GetBaseURL() string {
	return c.baseURL
}

func getEnvValue(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
