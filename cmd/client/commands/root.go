package commands

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/client/app"
	"github.com/sreway/gophkeeper/internal/config"
)

func NewRootCmd(ctx context.Context) *cobra.Command {
	var (
		cmd        *cobra.Command
		configFile string
		client     app.Client
	)
	cmd = &cobra.Command{
		Use:   "keeper-cli",
		Short: "gophkeeper",
		Long:  "Gophkeeper console password manager client.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.SetContext(ctx)
	cmd.PersistentFlags().StringVar(&configFile, "config", ".keeper_client.yaml", "Config file")
	_ = cmd.ParseFlags(os.Args[1:])

	clientConfig, err := config.NewClient(configFile)
	if err != nil {
		cobra.CheckErr(err)
	}

	client, err = app.NewClient(cmd.Context(), clientConfig)
	if err != nil {
		cobra.CheckErr(err)
	}

	cmd.AddCommand(newRegisterCmd(client))
	cmd.AddCommand(newLoginCmd(client))
	cmd.AddCommand(newLogoutCmd(client))
	cmd.AddCommand(newStatusCmd(client))
	cmd.AddCommand(newActionsCmd(client))
	cmd.AddCommand(newSecretCmd(client))
	return cmd
}
