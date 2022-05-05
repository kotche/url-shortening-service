package config

import "github.com/caarlos0/env/v6"

const ShortURLLen = 7

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL    string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FilePath   string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() (*Config, error) {
	conf := &Config{}
	err := env.Parse(conf)

	if err != nil {
		return nil, err
	}
	return conf, nil
}
