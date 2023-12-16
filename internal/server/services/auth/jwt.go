package auth

import (
	"time"

	"github.com/sreway/gophkeeper/internal/domain/models"
	"github.com/sreway/gophkeeper/internal/lib/jwt"
)

type (
	JWTManager interface {
		NewToken(user *models.User) (string, error)
		VerifyToken(token string) (*jwt.UserClaims, error)
	}

	jwtManager struct {
		secretKey string
		tokenTTL  time.Duration
	}
)

func (j *jwtManager) NewToken(user *models.User) (string, error) {
	return jwt.NewToken(user, j.secretKey, j.tokenTTL)
}

func (j *jwtManager) VerifyToken(token string) (*jwt.UserClaims, error) {
	return jwt.VerifyToken(j.secretKey, token)
}

func NewJWTManager(secretKey string, tokenTTL time.Duration) *jwtManager {
	return &jwtManager{
		secretKey: secretKey,
		tokenTTL:  tokenTTL,
	}
}
