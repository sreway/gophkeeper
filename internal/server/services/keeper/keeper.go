package keeper

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/sreway/gophkeeper/internal/domain/models"
	"github.com/sreway/gophkeeper/internal/server/services"
)

type (
	Service interface {
		HealthCheck(ctx context.Context) error
		GetServerID(ctx context.Context) (uuid.UUID, error)
		SetServerID(ctx context.Context, id uuid.UUID) error
		GetRecoveryKeyShare(ctx context.Context, userID uuid.UUID) ([]byte, error)
		UpdateSeal(ctx context.Context, userID uuid.UUID, seal *models.Seal) error
		ListUpdatedSecrets(ctx context.Context, ownerID uuid.UUID, updatedAfter *time.Time) ([]*models.Secret, error)
		BatchUpdateSecrets(ctx context.Context, userID uuid.UUID, secrets []*models.Secret) error
	}

	service struct {
		serverProvider services.ServerProvider
		sealProvider   services.SealProvider
		secretProvider services.SecretProvider
	}

	HealthChecker interface {
		HealthCheck(ctx context.Context) error
	}
)

func (s *service) HealthCheck(ctx context.Context) error {
	return s.serverProvider.HealthCheck(ctx)
}

func (s *service) GetServerID(ctx context.Context) (uuid.UUID, error) {
	return s.serverProvider.GetServerID(ctx)
}

func (s *service) SetServerID(ctx context.Context, id uuid.UUID) error {
	return s.serverProvider.SetServerID(ctx, id)
}

func (s *service) GetRecoveryKeyShare(ctx context.Context, userID uuid.UUID) ([]byte, error) {
	seal, err := s.sealProvider.GetSeal(ctx, userID)
	if err != nil {
		return nil, err
	}
	return seal.RecoveryShare, nil
}

func (s *service) UpdateSeal(ctx context.Context, userID uuid.UUID, seal *models.Seal) error {
	return s.sealProvider.UpdateSeal(ctx, userID, seal)
}

func (s *service) ListUpdatedSecrets(ctx context.Context, userID uuid.UUID, t *time.Time) ([]*models.Secret, error) {
	return s.secretProvider.ListUpdatedSecrets(ctx, userID, t)
}

func (s *service) BatchUpdateSecrets(ctx context.Context, userID uuid.UUID, secrets []*models.Secret) error {
	return s.secretProvider.BatchUpdateSecrets(ctx, userID, secrets)
}

func New(storage services.KeeperProvider) *service {
	return &service{
		serverProvider: storage,
		sealProvider:   storage,
		secretProvider: storage,
	}
}
