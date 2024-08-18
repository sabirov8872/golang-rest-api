package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	ServerAddress string
	DatabasePath  string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := Config{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		DatabasePath:  os.Getenv("DATABASE_PATH"),
	}

	return &config, nil
}
