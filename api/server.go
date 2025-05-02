package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"simplebank/config"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type (
	Server struct {
		webServer   *echo.Echo
		tasksClient *asynq.Client
		tasksServer *asynq.Server
	}
)

func NewServer(settings *config.Settings) *Server {
	db, err := sql.Open("mysql", settings.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	tasksClient := asynq.NewClient(asynq.RedisClientOpt{Addr: settings.RedisURL})

	return &Server{
		webServer:   NewWebServer(db, tasksClient),
		tasksClient: tasksClient,
		tasksServer: NewTasksServer(settings.RedisURL),
	}
}

func (s *Server) Start(addr string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := s.webServer.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Msg("shutting down web server")
		}
	}()
	log.Info().Msgf("web server running on %s", addr)

	go func() {
		mux := TasksServerMux()
		if err := s.tasksServer.Start(mux); err != nil {
			log.Fatal().Err(err).Msg("could not run tasks server")
		}
	}()
	log.Info().Msg("tasks server started processing")

	<-ctx.Done()

	log.Info().Msg("shutting down web server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.webServer.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Info().Msg("shutting down tasks server")
	s.tasksClient.Close()
	s.tasksServer.Shutdown()
}
