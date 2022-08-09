package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL    string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FilePath   string `env:"FILE_STORAGE_PATH"`
	DBConnect  string `env:"DATABASE_DSN"`
}

func NewConfig() (*Config, error) {
	conf := &Config{}

	if err := env.Parse(conf); err != nil {
		return nil, err
	}

	regStringVar(&conf.ServerAddr, "a", conf.ServerAddr, "server address")
	regStringVar(&conf.BaseURL, "b", conf.BaseURL, "base url")
	regStringVar(&conf.FilePath, "f", conf.FilePath, "file storage path")
	regStringVar(&conf.DBConnect, "d", conf.DBConnect, "database connection")
	flag.Parse()

	return conf, nil
}

func regStringVar(p *string, name string, value string, usage string) {
	if flag.Lookup(name) == nil {
		flag.StringVar(p, name, value, usage)
	}
}
