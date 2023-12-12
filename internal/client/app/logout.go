package app

import (
	"context"
)

func (c *client) Logout(ctx context.Context) error {
	profile, err := c.keeper.GetActiveProfile(ctx)
	if err != nil {
		return err
	}

	err = c.keeper.RemoveActiveProfile(ctx, c.id)
	if err != nil {
		return err
	}

	if profile.Session != nil {
		err = c.keeper.RemoveSession(ctx, profile.Session.ID)
		if err != nil {
			return err
		}
	}

	return c.keeper.RemoveActiveProfile(ctx, c.id)
}
