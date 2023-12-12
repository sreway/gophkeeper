package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	HashPassword string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
