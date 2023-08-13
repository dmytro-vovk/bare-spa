package config

import (
	"encoding/json"
	"os"

	"github.com/dmytro-vovk/envset"
)

type Config struct {
	WebServer WebServer `json:"webserver"`
	Database  Database  `json:"database"`
	Merchant  Merchant  `json:"merchant"`
}

type WebServer struct {
	Domain string `json:"domain"`
	TLS    TLS    `json:"tls"`
}

type TLS struct {
	Enabled   bool     `json:"enabled"`
	HostNames []string `json:"host_names"`
	CertDir   string   `json:"cert_dir"`
}

type Database struct {
	Host     string `env:"MYSQL_HOST"`
	Name     string `env:"MYSQL_NAME"`
	User     string `env:"MYSQL_USER"`
	Password string `env:"MYSQL_PASS"`
}

type Merchant struct {
	ID        int    `json:"id"`
	Key       string `env:"MERCHANT_KEY"`
	SystemKey string `env:"SYSTEM_KEY"`
}

func Load(fileName string) (*Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err := envset.Set(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
