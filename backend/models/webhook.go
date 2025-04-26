package models

import (
	"time"

	"github.com/google/uuid"
)

type Webhook struct {
	ID             uuid.UUID `json:"id"`
	EventType      string    `json:"event_type"`
	Payload        []byte    `json:"payload"`  // []byte because it's JSONB
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Status         string    `json:"status"`    // 'pending', 'delivered', 'failed'
	Delivered      bool      `json:"delivered"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Retries        int       `json:"retries"`
}