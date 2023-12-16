package auth

import (
	"context"

	"github.com/google/uuid"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
	"github.com/sreway/gophkeeper/internal/domain/models"
	"github.com/sreway/gophkeeper/internal/lib/crypt"
)

func (a *authGRPC) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	emailAddress, err := models.NewEmail(in.Email)
	if err != nil {
		return nil, err
	}

	if len(in.Password) == 0 {
		return nil, models.ErrInvalidPassword
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        emailAddress.String(),
		HashPassword: crypt.HashSum([]byte(in.Password), nil),
	}

	sealID, err := uuid.Parse(in.Seal.Id)
	if err != nil {
		return nil, err
	}

	seal := &models.Seal{
		ID:                 sealID,
		EncryptedShares:    in.Seal.EncryptedShares,
		RecoveryShare:      in.Seal.RecoveryShare,
		TotalShares:        in.Seal.TotalShares,
		RequiredShares:     in.Seal.RequiredShares,
		HashMasterPassword: in.Seal.HashMasterPassword,
		HashKey:            in.Seal.HashKey,
	}

	err = a.auth.Register(ctx, user, seal)
	if err != nil {
		return nil, err
	}

	response := &pb.RegisterResponse{
		User: &pb.User{
			Id:           user.ID.String(),
			Email:        user.Email,
			HashPassword: user.HashPassword,
		},
	}

	return response, nil
}
