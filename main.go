package main

import (
	"database/sql"
	"flag"
	"os"

	"simplebank/api"
	"simplebank/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	settings := config.New()

	db, err := sql.Open("mysql", settings.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	server := api.NewServer(db)

	err = server.Start(":8080")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
