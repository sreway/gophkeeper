package keeper

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
)

func (k *keeperGRPC) GetRecoveryKeyShare(ctx context.Context, in *pb.RecoveryKeyShareRequest) (*pb.RecoveryKeyShareResponse, error) {
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

	recoveryKeyShare, err := k.keeper.GetRecoveryKeyShare(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &pb.RecoveryKeyShareResponse{
		RecoveryKeyShare: recoveryKeyShare,
	}, nil
}
