package api

import (
	"context"
	"fmt"

	"simplebank/tasks"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewTasksServer(redisURL string) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisURL},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				tasks.QueueCritical: 10,
				tasks.QueueDefault:  5,
				tasks.QueueLow:      1,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().
					Err(err).
					Str("type", task.Type()).
					Bytes("payload", task.Payload()).
					Msg("task failed")
			}),
			Logger:   &Logger{},
			LogLevel: asynq.WarnLevel,
		},
	)
}

func TasksServerMux() *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)

	return mux
}

type Logger struct{}

func (l *Logger) Print(level zerolog.Level, args ...any) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (l *Logger) Debug(args ...any) {
	l.Print(zerolog.DebugLevel, args...)
}

func (l *Logger) Info(args ...any) {
	l.Print(zerolog.InfoLevel, args...)
}

func (l *Logger) Warn(args ...any) {
	l.Print(zerolog.WarnLevel, args...)
}

func (l *Logger) Error(args ...any) {
	l.Print(zerolog.ErrorLevel, args...)
}

func (l *Logger) Fatal(args ...any) {
	l.Print(zerolog.FatalLevel, args...)
}
