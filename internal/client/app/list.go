package app

import (
	"context"
	"errors"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

func (c *client) ListSecret(ctx context.Context) ([]models.Entry, error) {
	profile, err := c.keeper.GetActiveProfile(ctx)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return nil, models.ErrLoginRequired
		default:
			return nil, err
		}
	}

	secrets, err := c.keeper.ListSecret(ctx, profile.User.ID)
	if err != nil {
		return nil, err
	}

	entries := make([]models.Entry, 0, len(secrets))

	for _, secret := range secrets {
		var entry models.Entry
		entry, err = c.secretToEntry(secret)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
