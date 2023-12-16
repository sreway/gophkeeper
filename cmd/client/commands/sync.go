package commands

import (
	"errors"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/client/app"
	"github.com/sreway/gophkeeper/internal/domain/models"
)

func newSyncCmd(client app.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync secret database.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := client.ServerHealthCheck(cmd.Context())
			if err != nil {
				pterm.Error.Println("The operation cannot be performed because the server is unavailable.")
				return err
			}
			err = client.LoadActiveProfile(cmd.Context())
			if err != nil && !errors.Is(err, models.ErrInvalidMasterPassword) {
				pterm.Error.Println("Failed to log out of user account.")
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			resolveCH := make(chan int)
			conflictCH := make(chan []models.Entry)

			go func() {
				for {
					select {
					case entries := <-conflictCH:
						selected, err := viewConflictResolve(entries)
						if err != nil {
							return
						}
						resolveCH <- selected
					case <-cmd.Context().Done():
						return
					}
				}
			}()
			err := client.Sync(cmd.Context(), conflictCH, resolveCH)
			if err != nil {
				pterm.Error.Println("Failed sync")
				return err
			}
			pterm.Success.Println("Success sync")
			return nil
		},
	}
	return cmd
}
