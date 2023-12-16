package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
	Email  string
}

func NewToken(user *models.User, secret string, duration time.Duration) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(duration)),
		},
		UserID: user.ID,
		Email:  user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(secret, accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrInvalidSigningMethod
			}
			return []byte(secret), nil
		})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, ErrInvalidTokenClaims
	}

	return claims, err
}
