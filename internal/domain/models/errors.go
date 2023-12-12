package models

import (
	"errors"
)

var (
	ErrInvalidEmailAddress   = errors.New("invalid email address")
	ErrNotFound              = errors.New("not found")
	ErrAlreadyExists         = errors.New("already exists")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrInvalidMasterPassword = errors.New("invalid master password")
	ErrLoginRequired         = errors.New("login required")
	ErrServiceUnavailable    = errors.New("service unavailable")
	ErrRegisterRequired      = errors.New("register required")
	ErrUnknownSecretType     = errors.New("unknown secret type")
	ErrInvalidID             = errors.New("invalid id")
	ErrIDDoNotMatch          = errors.New("id do not match")
)
