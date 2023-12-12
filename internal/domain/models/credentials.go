package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type (
	Credentials struct {
		EntryID   uuid.UUID `json:"entry_id"`
		EntryName string    `json:"entry_name"`
		Username  string    `json:"username"`
		Password  string    `json:"password"`
	}
)

func (c *Credentials) ID() uuid.UUID {
	return c.EntryID
}

func (c *Credentials) Type() EntryType {
	return CredentialsType
}

func (c *Credentials) Bytes() ([]byte, error) {
	return json.Marshal(c)
}

func NewCredentials(name, username, password string) *Credentials {
	return &Credentials{
		EntryID:   uuid.New(),
		EntryName: name,
		Username:  username,
		Password:  password,
	}
}
