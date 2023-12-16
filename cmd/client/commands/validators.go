package commands

import (
	"github.com/pterm/pterm"

	"github.com/sreway/gophkeeper/internal/client/app"
)

func validateMasterPasswordSet(c app.Client) error {
	if ok := c.IsMasterPasswordExists(); !ok {
		masterPassword, err := pterm.DefaultInteractiveTextInput.WithMask("*").
			Show("Enter master password")
		if err != nil {
			return err
		}
		c.SetMasterPassword(masterPassword)
	}
	return nil
}
