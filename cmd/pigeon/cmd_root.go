package main

import (
	"github.com/palomachain/sparrow/app"
	"github.com/spf13/cobra"
)

// flags
var (
	flagConfigPath     string
	configRequiredCmds []*cobra.Command
)

var (
	rootCmd = &cobra.Command{
		Use:          "sparrow",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			found := false
			for _, curr := range configRequiredCmds {
				if curr == cmd {
					found = true
					break
				}
			}
			if found {
				app.SetConfigPath(flagConfigPath)
			}
		},
	}
)

func configRequired(cmd *cobra.Command) {
	for _, curr := range configRequiredCmds {
		if curr == cmd {
			return
		}
	}
	configRequiredCmds = append(configRequiredCmds, cmd)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagConfigPath, "config", "c", "~/.sparrow/config.yaml", "path to the config file")
}
