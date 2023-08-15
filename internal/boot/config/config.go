package config

import (
	"encoding/json"
	"os"

	"github.com/dmytro-vovk/envset"
)

type Config struct {
	WebServer WebServer `json:"webserver"`
}

type WebServer struct {
	Domain string `json:"domain" env:"LISTEN"`
	TLS    TLS    `json:"tls"`
}

type TLS struct {
	Enabled   bool     `json:"enabled"`
	HostNames []string `json:"host_names"`
	CertDir   string   `json:"cert_dir"`
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
