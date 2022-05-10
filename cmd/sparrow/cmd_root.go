package main

import (
	"github.com/palomachain/sparrow/app"
	"github.com/spf13/cobra"
)

// flags
var (
	flagConfigPath       string
	noConfigRequiredCmds []*cobra.Command
)

var (
	rootCmd = &cobra.Command{
		Use:          "sparrow",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			for _, curr := range noConfigRequiredCmds {
				if curr == cmd {
					return
				}
			}
			app.SetConfigPath(flagConfigPath)
		},
	}
)

func noConfigRequired(cmd *cobra.Command) {
	for _, curr := range noConfigRequiredCmds {
		if curr == cmd {
			return
		}
	}
	noConfigRequiredCmds = append(noConfigRequiredCmds, cmd)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagConfigPath, "config", "c", "~/.sparrow/config.yaml", "path to the config file")
}
