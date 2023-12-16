package commands

import (
	"errors"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/client/app"
	"github.com/sreway/gophkeeper/internal/domain/models"
)

func newSecretCmd(client app.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Secret management actions.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := client.LoadActiveProfile(cmd.Context())
			if err != nil && !errors.Is(err, models.ErrInvalidMasterPassword) {
				pterm.Error.Println("Failed to login of user account.")
				return err
			}
			return validateMasterPasswordSet(client)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newSecretCreateCmd(client))
	cmd.AddCommand(newSecretListCmd(client))
	cmd.AddCommand(newSecretShowCmd(client))
	cmd.AddCommand(newSecretUpdateCmd(client))
	cmd.AddCommand(newSecretDeleteCmd(client))
	return cmd
}
