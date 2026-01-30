package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type UserSub struct {
	ID           uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServiceName  string     `json:"service_name" example:"Netflix"`
	ServicePrice int        `json:"service_price" example:"990"`
	UserID       uuid.UUID  `json:"user_id" example:"550e8400-e29b-41d4-a716-446655441111"`
	StartedAt    time.Time  `json:"started_at" example:"2025-07-01T00:00:00Z"`
	EndedAt      *time.Time `json:"ended_at,omitempty" example:"2026-07-01T00:00:00Z"`
}
