package keeper

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (k *keeperGRPC) UpdateSeal(ctx context.Context, in *pb.UpdateSealRequest) (*pb.UpdateSealResponse, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	if len(md.Get("userID")) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "userID is not set")
	}

	userID, err := uuid.Parse(md.Get("userID")[0])
	if err != nil {
		return nil, err
	}

	seal := &models.Seal{
		EncryptedShares:    in.Seal.EncryptedShares,
		RecoveryShare:      in.Seal.RecoveryShare,
		TotalShares:        in.Seal.TotalShares,
		RequiredShares:     in.Seal.RequiredShares,
		HashMasterPassword: in.Seal.HashMasterPassword,
		HashKey:            in.Seal.HashKey,
	}

	return new(pb.UpdateSealResponse), k.keeper.UpdateSeal(ctx, userID, seal)
}
