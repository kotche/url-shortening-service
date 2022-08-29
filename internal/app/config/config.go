package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/caarlos0/env/v6"
)

// Config sets the basic settings
type (
	Config struct {
		ServerAddr  string `env:"SERVER_ADDRESS" envDefault:"localhost:8080" json:"server_address"`
		BaseURL     string `env:"BASE_URL" envDefault:"http://localhost:8080" json:"base_url"`
		FilePath    string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
		DBConnect   string `env:"DATABASE_DSN" json:"database_dsn"`
		EnableHTTPS bool   `env:"ENABLE_HTTPS" json:"enable_https"`

		HostWhitelist []string `json:"hostWhitelist"`
	}
)

//NewConfig priority: env and flag are on the same level, the configuration file is below
//example config file flag: -c ./internal/app/config/config.json
func NewConfig() (*Config, error) {
	conf := &Config{}

	flagServerAddr := flag.String("a", "", "SERVER_ADDRESS")
	flagBaseURL := flag.String("b", "", "BASE_URL")
	flagFilePath := flag.String("f", "", "FILE_STORAGE_PATH")
	flagDBConnect := flag.String("d", "", "DATABASE_DSN")
	flagEnableHTTPS := flag.String("s", "", "ENABLE_HTTPS")
	flagConfigFilePath := flag.String("c", "", "CONFIG")
	flagConfigFilePathTwo := flag.String("config", "", "CONFIG")
	flag.Parse()

	var configFilePath string
	if *flagConfigFilePath != "" {
		configFilePath = *flagConfigFilePath
	} else if *flagConfigFilePathTwo != "" {
		configFilePath = *flagConfigFilePathTwo
	}

	if configFilePath != "" {
		readConfigFile(conf, configFilePath)
	}

	if err := env.Parse(conf); err != nil {
		return nil, err
	}

	if *flagServerAddr != "" {
		conf.ServerAddr = *flagServerAddr
	}
	if *flagBaseURL != "" {
		conf.BaseURL = *flagBaseURL
	}
	if *flagFilePath != "" {
		conf.FilePath = *flagFilePath
	}
	if *flagDBConnect != "" {
		conf.DBConnect = *flagDBConnect
	}
	if *flagEnableHTTPS != "" {
		enableHTTPS, err := strconv.ParseBool(*flagEnableHTTPS)
		if err != nil {
			log.Printf("parse flagEnableHTTPS error: %s", err)
		} else {
			conf.EnableHTTPS = enableHTTPS
		}
	}

	return conf, nil
}

func readConfigFile(conf *Config, configFilePath string) {
	configFile, err := os.Open(configFilePath)
	if err != nil {
		log.Printf("file config open error: %s", err)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&conf)
	if err != nil {
		log.Printf("file config decode error: %s", err)
	}
}
