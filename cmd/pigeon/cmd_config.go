package main

import (
	"github.com/spf13/cobra"
	"github.com/palomachain/sparrow/app"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "config related commands",
	}
	validateConfigCmd = &cobra.Command{
		Use:   "validate",
		Short: "validates configuration",
		RunE: func(cmd *cobra.Command, args []string) error {

			_, err := app.PalomaClient().Keyring().List()
			if err != nil {
				return err
			}
			// TODO: add more checks!

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(validateConfigCmd)
}
