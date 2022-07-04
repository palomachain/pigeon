package main

import (
	"fmt"

	"github.com/palomachain/pigeon/app"
	"github.com/spf13/cobra"
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "starts the pigeon server",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("*chirp chirp*")
			return app.Relayer().Start(cmd.Context())
		},
	}
)

func init() {
	rootCmd.AddCommand(startCmd)

	configRequired(startCmd)
}
