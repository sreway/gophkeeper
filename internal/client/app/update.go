package app

import (
	"context"
	"errors"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (c *client) UpdateSecret(ctx context.Context, entry models.Entry) error {
	profile, err := c.keeper.GetActiveProfile(ctx)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return models.ErrLoginRequired
		default:
			return err
		}
	}

	secret, err := c.entryToSecret(profile.ID, entry)
	if err != nil {
		return err
	}
	err = c.keeper.UpdateSecret(ctx, secret)
	if err != nil {
		return err
	}

	return nil
}
