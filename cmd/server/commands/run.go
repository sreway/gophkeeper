package commands

import (
	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/server/app"
)

func newRunCmd(server app.Server) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run password manager server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.Run(cmd.Context())
		},
	}
	return cmd
}
