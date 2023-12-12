package commands

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/client/app"
)

func newSecretCreateCmd(client app.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create secret.",
		RunE: func(cmd *cobra.Command, args []string) error {
			entryType := viewSelectEntryType()
			entry, err := viewCreateEntry(entryType)
			if err != nil {
				return err
			}
			_, err = client.CreateSecret(cmd.Context(), entry)
			if err != nil {
				return err
			}

			pterm.Success.Println(fmt.Sprintf("Success create secret %s.", entry.ID().String()))
			return nil
		},
	}

	return cmd
}
