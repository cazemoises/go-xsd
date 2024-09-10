package config

import (
	"log"
	"os"
)

type Config struct {
	Port string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Port set to %s", port)
	return &Config{
		Port: port,
	}
}
