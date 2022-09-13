package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

// Config sets the basic settings
type Config struct {
	ServerAddr    string   `env:"SERVER_ADDRESS" envDefault:"localhost:8080" json:"server_address"`
	BaseURL       string   `env:"BASE_URL" envDefault:"http://localhost:8080" json:"base_url"`
	FilePath      string   `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DBConnect     string   `env:"DATABASE_DSN" json:"database_dsn"`
	EnableHTTPS   bool     `env:"ENABLE_HTTPS" json:"enable_https"`
	TrustedSubnet string   `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	HostWhitelist []string `json:"hostWhitelist"`
}

// NewConfig priority: env and flag are on the same level, the configuration file is below
// example config file flag: -c ./internal/app/config/config.json
func NewConfig() (*Config, error) {
	var serverAddr, baseURL, filePath, dbConnect, enableHTTPSStr, configFilePath, trustedSubnet string

	regStringVar(&serverAddr, "a", serverAddr, "server address")
	regStringVar(&baseURL, "b", baseURL, "base url")
	regStringVar(&filePath, "f", filePath, "file storage path")
	regStringVar(&dbConnect, "d", dbConnect, "database connection")
	regStringVar(&enableHTTPSStr, "s", enableHTTPSStr, "enable HTTPS")
	regStringVar(&configFilePath, "c", configFilePath, "config file")
	regStringVar(&configFilePath, "config", configFilePath, "config file")
	regStringVar(&trustedSubnet, "t", trustedSubnet, "trusted subnet")
	flag.Parse()

	conf := &Config{}

	if configFilePath != "" {
		readConfigFile(conf, configFilePath)
	}

	if err := env.Parse(conf); err != nil {
		return nil, err
	}

	if serverAddr != "" {
		conf.ServerAddr = serverAddr
	}
	if baseURL != "" {
		conf.BaseURL = baseURL
	}
	if filePath != "" {
		conf.FilePath = filePath
	}
	if dbConnect != "" {
		conf.DBConnect = dbConnect
	}
	if enableHTTPSStr == "true" {
		conf.EnableHTTPS = true
	}
	if trustedSubnet != "" {
		conf.TrustedSubnet = trustedSubnet
	}

	return conf, nil
}

func regStringVar(p *string, name string, value string, usage string) {
	if flag.Lookup(name) == nil {
		flag.StringVar(p, name, value, usage)
	}
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
