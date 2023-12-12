package commands

import (
	"errors"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/client/app"
	"github.com/sreway/gophkeeper/internal/domain/models"
)

func newLogoutCmd(client app.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of the current user account.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := client.LoadActiveProfile(cmd.Context())
			if err != nil && !errors.Is(err, models.ErrInvalidMasterPassword) {
				pterm.Error.Println("Failed to log out of user account.")
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := client.Logout(cmd.Context())
			if err != nil {
				pterm.Error.Println("Failed to log out of user account.")
				return err
			}

			pterm.Success.Println("Success log out a user account.")
			return nil
		},
	}
	return cmd
}
