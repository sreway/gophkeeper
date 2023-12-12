package commands

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/config"
	"github.com/sreway/gophkeeper/internal/server/app"
)

func NewRootCmd(ctx context.Context) *cobra.Command {
	var (
		cmd        *cobra.Command
		configFile string
	)

	cmd = &cobra.Command{
		Use:   "keeper-server",
		Short: "gophkeeper",
		Long:  "Gophkeeper password manager server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVar(&configFile, "config", ".keeper_server.yaml",
		"Config file")

	_ = cmd.ParseFlags(os.Args[1:])

	serverConfig, err := config.NewServer(configFile)
	if err != nil {
		cobra.CheckErr(err)
	}

	server, err := app.NewServer(ctx, serverConfig)
	if err != nil {
		cobra.CheckErr(err)
	}

	cmd.AddCommand(newRunCmd(server))
	return cmd
}
