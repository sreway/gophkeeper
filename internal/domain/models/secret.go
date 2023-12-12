package models

import (
	"time"

	"github.com/google/uuid"
)

type Secret struct {
	ID             uuid.UUID
	Owner          uuid.UUID
	EncryptedValue []byte
	Hash           string
	Type           EntryType
	IsDeleted      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
