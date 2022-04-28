package main

import (
	"github.com/spf13/cobra"
	"github.com/palomachain/sparrow/app"
)

// flags
var (
	flagConfigPath string
)

var (
	rootCmd = &cobra.Command{
		Use:          "sparrow",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			app.SetConfigPath(flagConfigPath)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagConfigPath, "config", "c", "~/.sparrow/config.yaml", "path to the config file")
}
