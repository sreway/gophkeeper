package commands

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/sreway/gophkeeper/internal/client/app"
)

func newLoginCmd(client app.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log into a user account.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
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

			_, err = client.Login(cmd.Context(), emailFlag, passwordFlag)
			if err != nil {
				pterm.Error.Println("Failed log into a user account.")
				return err
			}
			pterm.Success.Println("Success log into a user account.")

			return nil
		},
	}

	cmd.PersistentFlags().String("email", "", "Email address.")
	cmd.PersistentFlags().String("password", "", "Password.")

	return cmd
}
