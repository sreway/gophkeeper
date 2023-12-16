package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type (
	Text struct {
		EntryID   uuid.UUID `json:"entry_id"`
		EntryName string    `json:"entry_name"`
		Value     string    `json:"value"`
	}
)

func (t *Text) ID() uuid.UUID {
	return t.EntryID
}

func (t *Text) Type() EntryType {
	return TextType
}

func (t *Text) Bytes() ([]byte, error) {
	return json.Marshal(t)
}

func NewText(name, value string) *Text {
	return &Text{
		EntryID:   uuid.New(),
		EntryName: name,
		Value:     value,
	}
}
