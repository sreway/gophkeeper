package app

import (
	"context"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/sreway/gophkeeper/gen/go/gophkeeper/v1"
	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (c *client) Sync(ctx context.Context, conflictCH chan []models.Entry, resolveCH chan int) error {
	var lastSync *time.Time
	profile, err := c.keeper.GetActiveProfile(ctx)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return models.ErrLoginRequired
		default:
			return err
		}
	}

	lastSync, err = c.keeper.GetLastSync(ctx, profile.User.ID)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return err
	}

	localUpdates, err := c.keeper.ListUpdatedSecrets(ctx, profile.User.ID, lastSync)
	if err != nil {
		return err
	}

	updateAfter := new(timestamp.Timestamp)

	if lastSync != nil {
		updateAfter.Seconds = lastSync.Unix()
		updateAfter.Nanos = int32(lastSync.Nanosecond())
	}

	listUpdatedSecretsRequest := &pb.ListUpdatedSecretsRequest{
		UpdatedAfter: updateAfter,
	}

	listUpdatedSecretsResponse, err := c.keeperGRPC.ListUpdatedSecrets(ctx, listUpdatedSecretsRequest)
	if err != nil {
		return err
	}

	serverUpdates := make([]*models.Secret, len(listUpdatedSecretsResponse.Secret))
	conflicts := make([][]models.Entry, 0)
	for idx, pbSecret := range listUpdatedSecretsResponse.Secret {
		var id, owner uuid.UUID
		id, err = uuid.Parse(pbSecret.Id)
		if err != nil {
			return err
		}

		owner, err = uuid.Parse(pbSecret.Owner)
		if err != nil {
			return err
		}

		if owner.String() != profile.User.ID.String() {
			return models.ErrIDDoNotMatch
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

		serverUpdates[idx] = secret
	}

	for _, l := range localUpdates {
		for _, s := range serverUpdates {
			if l.ID.String() == s.ID.String() {
				var localEntry, remoteEntry models.Entry
				localEntry, err = c.secretToEntry(l)
				if err != nil {
					return err
				}

				remoteEntry, err = c.secretToEntry(s)
				if err != nil {
					return err
				}

				conflicts = append(conflicts, []models.Entry{
					localEntry, remoteEntry,
				})
			}
		}
	}

	for _, conflict := range conflicts {
		conflictCH <- conflict
		resolveID := <-resolveCH
		switch resolveID {
		case 0:
			for idx, secret := range serverUpdates {
				if secret.ID == conflict[resolveID].ID() {
					serverUpdates = append(serverUpdates[:idx], serverUpdates[idx+1:]...)
					break
				}
			}
		case 1:
			for idx, secret := range localUpdates {
				if secret.ID == conflict[resolveID].ID() {
					localUpdates = append(localUpdates[:idx], localUpdates[idx+1:]...)
					break
				}
			}
		}
	}

	for _, secret := range localUpdates {
		if secret.IsDeleted {
			serverUpdates = append(serverUpdates, secret) // nozero
		}
	}

	batchUpdateSecretsRequest := new(pb.BatchUpdateSecretsRequest)

	batchUpdateSecretsRequest.Secret = make([]*pb.Secret, len(localUpdates))

	for idx, secret := range localUpdates {
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
		batchUpdateSecretsRequest.Secret[idx] = pbSecret
	}

	_, err = c.keeperGRPC.BatchUpdateSecrets(ctx, batchUpdateSecretsRequest)
	if err != nil {
		return err
	}

	err = c.keeper.BatchUpdateSecrets(ctx, profile.User.ID, serverUpdates)
	if err != nil {
		return err
	}

	_, err = c.keeper.UpdateLastSync(ctx, profile.User.ID)
	if err != nil {
		return err
	}

	return nil
}
