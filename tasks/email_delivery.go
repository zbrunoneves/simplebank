package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type EmailDeliveryPayload struct {
	Email string
}

func NewEmailDeliveryTask(email string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{Email: email})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailDelivery, payload), nil
}

func HandleEmailDeliveryTask(_ context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Info().Msgf("sending email to %s", p.Email)

	// Email delivery code ...

	return nil
}
