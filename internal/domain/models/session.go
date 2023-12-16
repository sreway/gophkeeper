package models

import (
	"github.com/google/uuid"
)

type Session struct {
	ID             uuid.UUID
	EncryptedToken []byte
}
