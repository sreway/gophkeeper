package auth

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (a *authGRPC) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	emailAddress, err := models.NewEmail(in.Email)
	if err != nil {
		return nil, err
	}

	if len(in.Password) == 0 {
		return nil, models.ErrInvalidPassword
	}

	user, seal, err := a.auth.Login(ctx, emailAddress.String(), in.Password)
	if err != nil && errors.Is(err, models.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, models.ErrNotFound.Error())
	}
	if err != nil {
		return nil, err
	}

	token, err := a.jwt.NewToken(user)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token: token,
		User: &pb.User{
			Id:           user.ID.String(),
			Email:        user.Email,
			HashPassword: user.HashPassword,
		},
		Seal: &pb.Seal{
			Id:                 seal.ID.String(),
			EncryptedShares:    seal.EncryptedShares,
			RecoveryShare:      seal.RecoveryShare,
			TotalShares:        seal.TotalShares,
			RequiredShares:     seal.RequiredShares,
			HashMasterPassword: seal.HashMasterPassword,
			HashKey:            seal.HashKey,
		},
	}, nil
}
