package models

import (
	"github.com/google/uuid"
)

type (
	EntryType     uint8
	EntryConflict struct {
		Local  Entry
		Server Entry
	}
)

const (
	CredentialsType EntryType = iota + 1
	TextType
)

func (et EntryType) String() string {
	switch et {
	case CredentialsType:
		return "Credentials"
	case TextType:
		return "Text"
	default:
		return "Unknown"
	}
}

type Entry interface {
	ID() uuid.UUID
	Bytes() ([]byte, error)
	Type() EntryType
}
