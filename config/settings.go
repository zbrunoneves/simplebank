package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Settings struct {
	DatabaseURL string
}

func New() *Settings {
	err := godotenv.Load()
	if err != nil {
		log.Warn().Msg("error parsing .env file")
	}

	return &Settings{
		DatabaseURL: os.Getenv("MYSQL_URL"),
	}
}
