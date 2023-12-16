package jwt

import (
	"errors"
)

var (
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidTokenClaims   = errors.New("invalid token claims")
)
