package configs

import (
	"os"
)

type Config struct {
	Port                string
	Environment         string
	OpenAIKey           string
	ProductSearchAPIKey string
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

	openAIKey := os.Getenv("OPENAI_API_KEY")
	productSearchAPIKey := os.Getenv("PRODUCT_SEARCH_API_KEY")

	return &Config{
		Port:                port,
		Environment:         environment,
		OpenAIKey:           openAIKey,
		ProductSearchAPIKey: productSearchAPIKey,
	}
}
