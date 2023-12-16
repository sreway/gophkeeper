package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/client/app"
)

func getStatus(ctx context.Context, c app.Client) string {
	var (
		serverStatus string
		activeUser   string
	)
	err := c.ServerHealthCheck(ctx)
	if err != nil {
		serverStatus = "Unavailable"
	} else {
		serverStatus = "Available"
	}

	currentUser, err := c.ActiveUser(ctx)
	if err != nil {
		activeUser = "No active user"
	} else {
		activeUser = currentUser.Email
	}

	statusItems := []string{
		fmt.Sprintf("ClientID: %s", c.ID().String()),
		fmt.Sprintf("Active User: %s", activeUser),
		fmt.Sprintf("Server: %s", serverStatus),
	}

	return strings.Join(statusItems, "\n")
}

func newStatusCmd(client app.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Get password manager status.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := validateMasterPasswordSet(client)
			if err != nil {
				return err
			}

			err = client.LoadActiveProfile(cmd.Context())
			if err != nil {
				pterm.Error.Println("Failed get password manager status.")
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			panel := pterm.DefaultBox.WithTitle("Status").Sprint(getStatus(cmd.Context(), client))
			return pterm.DefaultPanel.WithPanels([][]pterm.Panel{{{Data: panel}}}).Render()
		},
	}
	return cmd
}
