package app

import (
	"context"

	"google.golang.org/grpc/codes"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (c *client) ResetMasterPassword(ctx context.Context, email, password, masterPassword string) error {
	var (
		profile   *models.Profile
		masterKey []byte
	)

	emailAddress, err := models.NewEmail(email)
	if err != nil {
		return err
	}

	loginRequest := &pb.LoginRequest{
		Email:    email,
		Password: password,
	}

	loginResponse, err := c.authGRPC.Login(ctx, loginRequest)
	if err != nil {
		switch grpcGetStatusCode(err) {
		case codes.NotFound.String():
			return models.ErrRegisterRequired
		case codes.Unavailable.String():
			return models.ErrServiceUnavailable
		default:
			return err
		}
	}

	profile, err = c.keeper.GetProfile(ctx, emailAddress.String())
	if err != nil {
		return models.ErrRegisterRequired
	}

	c.tokenProvider.SetToken(loginResponse.Token)

	recoveryKeyResponse, err := c.keeperGRPC.GetRecoveryKeyShare(ctx, new(pb.RecoveryKeyShareRequest))
	if err != nil {
		return err
	}

	recoveryShares := [][]byte{
		profile.Seal.RecoveryShare, recoveryKeyResponse.RecoveryKeyShare,
	}

	masterKey, err = c.keeper.RecoveryMasterKey(recoveryShares)
	if err != nil {
		return err
	}

	local, remote, err := c.keeper.CreateKeySeals(masterKey, masterPassword)
	if err != nil {
		return err
	}

	local.ID = profile.Seal.ID

	updateSealRequest := &pb.UpdateSealRequest{
		Seal: &pb.Seal{
			EncryptedShares:    remote.EncryptedShares,
			RecoveryShare:      remote.RecoveryShare,
			RequiredShares:     remote.RequiredShares,
			TotalShares:        remote.TotalShares,
			HashMasterPassword: remote.HashMasterPassword,
			HashKey:            remote.HashKey,
		},
	}

	_, err = c.keeperGRPC.UpdateSeal(ctx, updateSealRequest)
	if err != nil {
		return err
	}

	return c.keeper.UpdateSeal(ctx, local)
}
