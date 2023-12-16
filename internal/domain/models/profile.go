package models

import (
	"github.com/google/uuid"
)

type Profile struct {
	ID      uuid.UUID
	User    *User
	Seal    *Seal
	Session *Session
}
