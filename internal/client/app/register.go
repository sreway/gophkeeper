package app

import (
	"context"

	"github.com/google/uuid"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (c *client) Register(ctx context.Context, email, password string) error {
	emailAddress, err := models.NewEmail(email)
	if err != nil {
		return err
	}

	local, remote, err := c.keeper.CreateKeySeals(nil, c.config.MasterPassword)
	if err != nil {
		return err
	}

	registerRequest := &pb.RegisterRequest{
		Email:    emailAddress.String(),
		Password: password,
		Seal: &pb.Seal{
			Id:                 remote.ID.String(),
			EncryptedShares:    remote.EncryptedShares,
			RecoveryShare:      remote.RecoveryShare,
			RequiredShares:     remote.RequiredShares,
			TotalShares:        remote.TotalShares,
			HashMasterPassword: remote.HashMasterPassword,
			HashKey:            remote.HashKey,
		},
	}

	registerResponse, err := c.authGRPC.Register(ctx, registerRequest)
	if err != nil {
		return err
	}

	userID, err := uuid.Parse(registerResponse.User.Id)
	if err != nil {
		return err
	}

	user := &models.User{
		ID:           userID,
		Email:        registerResponse.User.Email,
		HashPassword: registerResponse.User.HashPassword,
	}

	_, err = c.keeper.CreateProfile(ctx, user, local)
	if err != nil {
		return err
	}

	return nil
}
