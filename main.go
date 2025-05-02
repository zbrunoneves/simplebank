package main

import (
	"flag"
	"os"

	"simplebank/api"
	"simplebank/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	server := api.NewServer(settings)
	server.Start(":8081")
	// no code bellow this line will execute until server shutdown
}
