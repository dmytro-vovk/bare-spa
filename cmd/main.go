package main

import (
	"os"

	boot "github.com/Sergii-Kirichok/DTekSpeachParser/internal/boot"
)

func main() {
	configFile := "config.json"
	if cfg := os.Getenv("CONFIG"); cfg != "" {
		configFile = cfg
	}
	c, shutdown := boot.New(configFile)
	defer shutdown()

	c.Logger()
	if err := c.Webserver().Serve(); err != nil {
		panic(err)
	}
}
