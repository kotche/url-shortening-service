package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

// Config sets the basic settings
type Config struct {
	ServerAddr  string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL     string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FilePath    string `env:"FILE_STORAGE_PATH"`
	DBConnect   string `env:"DATABASE_DSN"`
	EnableHTTPS bool   `env:"ENABLE_HTTPS"`

	HostWhitelist []string
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
	regBoolVar(&conf.EnableHTTPS, "s", conf.EnableHTTPS, "enable HTTPS")
	flag.Parse()

	conf.HostWhitelist = []string{
		"localhost:8080",
	}

	return conf, nil
}

func regStringVar(p *string, name string, value string, usage string) {
	if flag.Lookup(name) == nil {
		flag.StringVar(p, name, value, usage)
	}
}

func regBoolVar(p *bool, name string, value bool, usage string) {
	if flag.Lookup(name) == nil {
		flag.BoolVar(p, name, value, usage)
	}
}
