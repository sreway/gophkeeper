package app

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (c *client) GetSecret(ctx context.Context, id uuid.UUID) (models.Entry, error) {
	profile, err := c.keeper.GetActiveProfile(ctx)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return nil, models.ErrLoginRequired
		default:
			return nil, err
		}
	}

	secret, err := c.keeper.GetSecret(ctx, id, profile.User.ID)
	if err != nil {
		return nil, err
	}

	return c.secretToEntry(secret)
}
