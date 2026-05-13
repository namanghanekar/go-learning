package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	KafkaBroker string
}

func LoadConfig() *Config {

	err := godotenv.Load()

	if err != nil {
		log.Println(".env file not found")
	}

	return &Config{
		AppPort: os.Getenv("APP_PORT"),

		KafkaBroker: os.Getenv("KAFKA_BROKER"),
	}
}
