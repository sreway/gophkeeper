package app

import (
	"context"
	"errors"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (c *client) CreateSecret(ctx context.Context, entry models.Entry) (*models.Secret, error) {
	profile, err := c.keeper.GetActiveProfile(ctx)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return nil, models.ErrLoginRequired
		default:
			return nil, err
		}
	}

	secret, err := c.entryToSecret(profile.User.ID, entry)
	if err != nil {
		return nil, err
	}
	err = c.keeper.AddSecret(ctx, secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
