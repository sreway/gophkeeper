package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

type (
	KeeperProvider interface {
		ServerProvider
		UserProvider
		SealProvider
		SecretProvider
	}

	ServerProvider interface {
		GetServerID(ctx context.Context) (uuid.UUID, error)
		SetServerID(ctx context.Context, id uuid.UUID) error
		HealthCheck(ctx context.Context) error
	}

	UserProvider interface {
		AddUser(ctx context.Context, user *models.User) error
		FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	}

	SealProvider interface {
		AddSeal(ctx context.Context, userID uuid.UUID, seal *models.Seal) error
		GetSeal(ctx context.Context, userID uuid.UUID) (*models.Seal, error)
		UpdateSeal(ctx context.Context, userID uuid.UUID, seal *models.Seal) error
	}

	SecretProvider interface {
		ListUpdatedSecrets(ctx context.Context, ownerID uuid.UUID, updatedAfter *time.Time) ([]*models.Secret, error)
		BatchUpdateSecrets(ctx context.Context, userID uuid.UUID, secrets []*models.Secret) error
	}
)
