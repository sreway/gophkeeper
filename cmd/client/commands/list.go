package commands

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/client/app"
)

func newSecretListCmd(client app.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List secrets.",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := client.ListSecret(cmd.Context())
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				pterm.Warning.Println("Empty secrets.")
				return nil
			}
			return viewListEntry(entries)
		},
	}

	return cmd
}
