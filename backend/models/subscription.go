package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID `json:"id"`
	TargetURL   string    `json:"target_url"`
	Secret      string    `json:"secret,omitempty"`
	EventTypes  []string  `json:"event_types"` // List of event types this subscription is interested in
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
