package keeper

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (k *keeperGRPC) ListUpdatedSecrets(ctx context.Context, in *pb.ListUpdatedSecretsRequest) (*pb.ListUpdatedSecretsResponse, error) {
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

	updatedAfter := new(time.Time)

	if in.UpdatedAfter != nil {
		*updatedAfter = time.Unix(in.UpdatedAfter.Seconds, int64(in.UpdatedAfter.Nanos))
	}

	secrets, err := k.keeper.ListUpdatedSecrets(ctx, userID, updatedAfter)
	if err != nil {
		return nil, err
	}

	response := new(pb.ListUpdatedSecretsResponse)

	response.Secret = make([]*pb.Secret, len(secrets))

	for idx, secret := range secrets {
		pbSecret := &pb.Secret{
			Id:             secret.ID.String(),
			Owner:          secret.Owner.String(),
			EncryptedValue: secret.EncryptedValue,
			Hash:           secret.Hash,
			Type:           uint64(secret.Type),
			IsDeleted:      secret.IsDeleted,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: secret.CreatedAt.Unix(),
				Nanos:   int32(secret.CreatedAt.Nanosecond()),
			},
			UpdatedAt: &timestamppb.Timestamp{
				Seconds: secret.UpdatedAt.Unix(),
				Nanos:   int32(secret.UpdatedAt.Nanosecond()),
			},
		}
		response.Secret[idx] = pbSecret
	}

	return response, nil
}

func (k *keeperGRPC) BatchUpdateSecrets(ctx context.Context, in *pb.BatchUpdateSecretsRequest) (*pb.BatchUpdateSecretsResponse, error) {
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

	secrets := make([]*models.Secret, len(in.Secret))

	for idx, pbSecret := range in.Secret {
		var id, owner uuid.UUID
		id, err = uuid.Parse(pbSecret.Id)
		if err != nil {
			return nil, err
		}

		owner, err = uuid.Parse(pbSecret.Owner)
		if err != nil {
			return nil, err
		}

		if owner.String() != userID.String() {
			return nil, models.ErrIDDoNotMatch
		}

		secret := &models.Secret{
			ID:             id,
			Owner:          owner,
			EncryptedValue: pbSecret.EncryptedValue,
			Hash:           pbSecret.Hash,
			Type:           models.EntryType(pbSecret.Type),
			IsDeleted:      pbSecret.IsDeleted,
			CreatedAt:      time.Unix(pbSecret.CreatedAt.Seconds, int64(pbSecret.CreatedAt.Nanos)),
			UpdatedAt:      time.Unix(pbSecret.UpdatedAt.Seconds, int64(pbSecret.UpdatedAt.Nanos)),
		}

		secrets[idx] = secret
	}

	err = k.keeper.BatchUpdateSecrets(ctx, userID, secrets)
	if err != nil {
		return nil, err
	}

	return new(pb.BatchUpdateSecretsResponse), nil
}
