package models

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryLog struct {
	ID             int64 `json:"id"`
	WebhookID      uuid.UUID `json:"webhook_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	TargetURL      string    `json:"target_url"`
	AttemptNumber  int       `json:"attempt_number"`
	Status         string    `json:"status"`
	HTTPStatusCode int       `json:"http_status_code"`
	ErrorDetails   string    `json:"error_details"`
	Timestamp      time.Time `json:"timestamp"`
}


