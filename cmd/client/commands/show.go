package commands

import (
	"github.com/google/uuid"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/client/app"
	"github.com/sreway/gophkeeper/internal/domain/models"
)

func newSecretShowCmd(client app.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "show secret.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				secretIDFlag   string
				secretID       uuid.UUID
				jsonOutputFlag bool
				err            error
			)
			secretIDFlag, err = cmd.Flags().GetString("id")
			if err != nil {
				return err
			}

			if len(secretIDFlag) == 0 {
				secretIDFlag, _ = pterm.DefaultInteractiveTextInput.WithDefaultText("Secret ID").
					WithMultiLine(false).Show()
			}

			secretID, err = uuid.Parse(secretIDFlag)
			if err != nil {
				return models.ErrInvalidID
			}

			entry, err := client.GetSecret(cmd.Context(), secretID)
			if err != nil {
				return err
			}

			jsonOutputFlag, err = cmd.Flags().GetBool("json")
			if err != nil {
				return err
			}

			err = viewEntryOutput(entry, jsonOutputFlag)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.PersistentFlags().String("id", "", "Secret ID.")
	cmd.PersistentFlags().Bool("json", false, "JSON output.")
	return cmd
}
