package main

import (
	"log"
	"os"

	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/boot"
)

func main() {
	configPath := "config.json"
	if cfgPath := os.Getenv("CONFIG"); cfgPath != "" {
		configPath = cfgPath
	}

	c, err := boot.New(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Webserver().Serve("Web server"); err != nil {
		log.Fatal(err)
	}
}
