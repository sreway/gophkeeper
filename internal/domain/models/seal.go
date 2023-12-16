package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Seal struct {
		ID                 uuid.UUID
		EncryptedShares    []byte
		RecoveryShare      []byte
		TotalShares        uint64
		RequiredShares     uint64
		HashMasterPassword string
		HashKey            string
		CreatedAt          time.Time
		UpdatedAt          time.Time
	}
)
