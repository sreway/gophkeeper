package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
	"github.com/sreway/gophkeeper/internal/domain/models"
	"github.com/sreway/gophkeeper/internal/lib/crypt"
)

func (c *client) Login(ctx context.Context, email, password string) (*models.Profile, error) {
	var (
		serverUnavailable, userExist bool
		user                         *models.User
		seal                         *models.Seal
		profile                      *models.Profile
		masterKey, encryptedToken    []byte
	)

	emailAddress, err := models.NewEmail(email)
	if err != nil {
		return nil, err
	}

	hashPassword := crypt.HashSum([]byte(password), nil)

	loginRequest := &pb.LoginRequest{
		Email:    email,
		Password: password,
	}

	loginResponse, err := c.authGRPC.Login(ctx, loginRequest)
	if err != nil {
		switch grpcGetStatusCode(err) {
		case codes.NotFound.String():
			return nil, models.ErrRegisterRequired
		case codes.Unavailable.String():
			serverUnavailable = true
		default:
			return nil, err
		}
	}

	profile, err = c.keeper.GetProfile(ctx, emailAddress.String())
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}

	if err == nil {
		userExist = true
	}

	switch {
	case serverUnavailable && !userExist:
		return nil, models.ErrServiceUnavailable
	case !serverUnavailable && !userExist:
		var (
			userID, sealID uuid.UUID
			keyShares      [][]byte
		)

		userID, err = uuid.Parse(loginResponse.User.Id)
		if err != nil {
			return nil, err
		}

		user = &models.User{
			ID:           userID,
			Email:        loginResponse.User.Email,
			HashPassword: loginResponse.User.HashPassword,
		}

		sealID, err = uuid.Parse(loginResponse.Seal.Id)
		if err != nil {
			return nil, err
		}

		keyShares, err = c.keeper.DecryptKeyShares(loginResponse.Seal.EncryptedShares, c.config.MasterPassword)
		if err != nil {
			return nil, err
		}

		seal = &models.Seal{
			ID:                 sealID,
			EncryptedShares:    loginResponse.Seal.EncryptedShares,
			TotalShares:        loginResponse.Seal.TotalShares,
			RequiredShares:     loginResponse.Seal.RequiredShares,
			HashMasterPassword: loginResponse.Seal.HashMasterPassword,
			HashKey:            loginResponse.Seal.HashKey,
			RecoveryShare:      keyShares[0],
		}

		profile, err = c.keeper.CreateProfile(ctx, user, seal)
		if err != nil {
			return nil, err
		}
	}

	if profile.User.HashPassword != hashPassword {
		return nil, models.ErrInvalidPassword
	}

	if serverUnavailable {
		return nil, c.keeper.SetActiveProfile(ctx, c.id, profile.ID)
	}

	masterKey, err = c.keeper.DecryptMasterKey(profile.Seal, c.config.MasterPassword)
	if err != nil {
		return nil, err
	}

	encryptedToken, err = crypt.EncryptGCM([]byte(loginResponse.GetToken()), masterKey)
	if err != nil {
		return nil, err
	}

	err = c.keeper.CreateSession(ctx, profile.ID, encryptedToken)
	if err != nil {
		return nil, err
	}

	return nil, c.keeper.SetActiveProfile(ctx, c.id, profile.ID)
}
