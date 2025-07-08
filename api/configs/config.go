package configs

import (
	"os"
)

type Config struct {
	Port string
	Environment string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	return &Config{
		Port:        port,
		Environment: environment,
	}
}