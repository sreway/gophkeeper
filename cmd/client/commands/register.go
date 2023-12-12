package commands

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/client/app"
)

func newRegisterCmd(client app.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Create a user account.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := client.ServerHealthCheck(cmd.Context())
			if err != nil {
				pterm.Error.Println("The operation cannot be performed because the server is unavailable.")
				return err
			}
			return validateMasterPasswordSet(client)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				emailFlag    string
				passwordFlag string
				err          error
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

			err = client.Register(cmd.Context(), emailFlag, passwordFlag)
			if err != nil {
				pterm.Error.Println("Failed register user.")
				return err
			}

			pterm.Success.Println("Success register user.")

			return nil
		},
	}

	cmd.PersistentFlags().String("email", "", "Email address.")
	cmd.PersistentFlags().String("password", "", "Password.")
	return cmd
}
