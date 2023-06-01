package main

import (
	"fmt"

	"github.com/palomachain/pigeon/app"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints version info",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("App version:", app.Version())
		fmt.Println("Build commit hash:", app.Commit())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
