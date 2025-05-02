package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Settings struct {
	AppEnv      string
	DatabaseURL string
	RedisURL    string
}

func New() *Settings {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	err := godotenv.Load(fmt.Sprintf("%s.env", env))
	if err != nil {
		log.Warn().Err(err).Send()
	}

	return &Settings{
		AppEnv:      env,
		DatabaseURL: os.Getenv("MYSQL_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
	}
}
