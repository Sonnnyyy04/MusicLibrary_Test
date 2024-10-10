package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBUrl       string
	ExternalAPI string
	ServerPort  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Not found .env file")
	}
	return &Config{
		DBUrl:       os.Getenv("DATABASE_URL"),
		ExternalAPI: os.Getenv("EXTERNAL_API_URL"),
		ServerPort:  os.Getenv("SERVER_PORT"),
	}
}
