package commands

import (
	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/client/app"
)

func newActionsCmd(client app.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "actions",
		Short: "Additional actions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newResetCmd(client))
	cmd.AddCommand(newSyncCmd(client))
	return cmd
}
