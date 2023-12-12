package auth

import (
	"context"

	"github.com/sreway/gophkeeper/internal/domain/models"
	"github.com/sreway/gophkeeper/internal/lib/crypt"
	"github.com/sreway/gophkeeper/internal/server/services"
)

type (
	Service interface {
		Register(ctx context.Context, user *models.User, seal *models.Seal) error
		Login(ctx context.Context, email string, password string) (user *models.User, seal *models.Seal, err error)
	}

	service struct {
		userProvider services.UserProvider
		sealProvider services.SealProvider
	}
)

func (s *service) Register(ctx context.Context, user *models.User, seal *models.Seal) error {
	err := s.userProvider.AddUser(ctx, user)
	if err != nil {
		return err
	}

	err = s.sealProvider.AddSeal(ctx, user.ID, seal)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Login(ctx context.Context, email, password string) (user *models.User, seal *models.Seal, err error) {
	user, err = s.userProvider.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	if user.HashPassword != crypt.HashSum([]byte(password), nil) {
		return nil, nil, models.ErrInvalidPassword
	}

	seal, err = s.sealProvider.GetSeal(ctx, user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, seal, nil
}

func New(storage services.KeeperProvider) *service {
	return &service{
		userProvider: storage,
		sealProvider: storage,
	}
}
