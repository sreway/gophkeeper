package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

type (
	KeeperProvider interface {
		ClientProvider
		ProfileProvider
		SecretProvider
		SyncProvider
	}

	ClientProvider interface {
		GetClientID(ctx context.Context) (uuid.UUID, error)
		SetClientID(ctx context.Context, id uuid.UUID) error
	}

	ProfileProvider interface {
		AddProfile(ctx context.Context, profile *models.Profile) error
		GetProfile(ctx context.Context, email string) (*models.Profile, error)
		GetActiveProfile(ctx context.Context) (*models.Profile, error)
		SetActiveProfile(ctx context.Context, clientID, profileID uuid.UUID) error
		RemoveActiveProfile(ctx context.Context, clientID uuid.UUID) error
		AddSession(ctx context.Context, profileID uuid.UUID, session *models.Session) error
		RemoveSession(ctx context.Context, id uuid.UUID) error
		UpdateSeal(ctx context.Context, seal *models.Seal) error
	}

	SecretProvider interface {
		AddSecret(ctx context.Context, secret *models.Secret) error
		UpdateSecret(ctx context.Context, secret *models.Secret) error
		DeleteSecret(ctx context.Context, id, ownerID uuid.UUID) error
		GetSecret(ctx context.Context, id, ownerID uuid.UUID) (*models.Secret, error)
		ListSecret(ctx context.Context, ownerID uuid.UUID) ([]*models.Secret, error)
		ListUpdatedSecrets(ctx context.Context, userID uuid.UUID, updatedAfter *time.Time) ([]*models.Secret, error)
		BatchUpdateSecrets(ctx context.Context, userID uuid.UUID, secrets []*models.Secret) error
	}

	SyncProvider interface {
		GetLastSync(ctx context.Context, userID uuid.UUID) (*time.Time, error)
		UpdateLastSync(ctx context.Context, userID uuid.UUID) (*time.Time, error)
	}
)
