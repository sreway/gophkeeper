package commands

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/client/app"
	"github.com/sreway/gophkeeper/internal/domain/models"
)

func newSecretUpdateCmd(client app.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update secret.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				secretIDFlag string
				secretID     uuid.UUID
				err          error
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

			entry, err = viewUpdateEntry(entry)
			if err != nil {
				return err
			}

			err = client.UpdateSecret(cmd.Context(), entry)
			if err != nil {
				return err
			}
			pterm.Success.Println(fmt.Sprintf("Success update secret %s.", entry.ID().String()))
			return nil
		},
	}
	cmd.PersistentFlags().String("id", "", "Secret ID.")
	return cmd
}
