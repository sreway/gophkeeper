package commands

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/client/app"
)

func newResetMasterPassword(client app.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mp",
		Short: "Reset master password.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := client.ServerHealthCheck(cmd.Context())
			if err != nil {
				pterm.Error.Println("The operation cannot be performed because the server is unavailable.")
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				newMasterPassword string
				emailFlag         string
				passwordFlag      string
				err               error
			)

			emailFlag, err = cmd.Flags().GetString("email")
			if err != nil {
				return err
			}

			if len(emailFlag) == 0 {
				emailFlag, err = pterm.DefaultInteractiveTextInput.Show("Enter email address")
				if err != nil {
					return err
				}
			}

			passwordFlag, err = cmd.Flags().GetString("password")
			if err != nil {
				return err
			}

			if len(passwordFlag) == 0 {
				passwordFlag, err = pterm.DefaultInteractiveTextInput.WithMask("*").
					Show("Enter your password")
				if err != nil {
					return err
				}
			}

			newMasterPassword, err = pterm.DefaultInteractiveTextInput.WithMask("*").
				Show("Enter new master password")
			if err != nil {
				return err
			}

			err = client.ResetMasterPassword(cmd.Context(), emailFlag, passwordFlag, newMasterPassword)
			if err != nil {
				return err
			}
			pterm.Success.Println("Success reset master password.")
			return nil
		},
	}

	cmd.PersistentFlags().String("email", "", "Email address.")
	cmd.PersistentFlags().String("password", "", "Password.")

	return cmd
}

func newResetCmd(client app.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset actions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newResetMasterPassword(client))
	return cmd
}
