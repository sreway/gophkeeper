package app

import (
	"context"
	"errors"

	"github.com/sreway/gophkeeper/internal/domain/models"
	"github.com/sreway/gophkeeper/internal/lib/crypt"
)

func (c *client) LoadActiveProfile(ctx context.Context) error {
	profile, err := c.keeper.GetActiveProfile(ctx)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return models.ErrLoginRequired
		default:
			return err
		}
	}

	c.masterKey, err = c.keeper.DecryptMasterKey(profile.Seal, c.config.MasterPassword)
	if err != nil {
		return err
	}
	if profile.Session.EncryptedToken != nil {
		var token []byte
		token, err = crypt.DecryptGCM(profile.Session.EncryptedToken, c.masterKey)
		if err != nil {
			return err
		}
		c.tokenProvider.SetToken(string(token))
	}

	return nil
}
