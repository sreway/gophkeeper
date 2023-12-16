package app

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (c *client) DeleteSecret(ctx context.Context, id uuid.UUID) error {
	profile, err := c.keeper.GetActiveProfile(ctx)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return models.ErrLoginRequired
		default:
			return err
		}
	}

	secret, err := c.keeper.GetSecret(ctx, id, profile.User.ID)
	if err != nil {
		return err
	}

	secret.IsDeleted = true

	return c.keeper.UpdateSecret(ctx, secret)
}
