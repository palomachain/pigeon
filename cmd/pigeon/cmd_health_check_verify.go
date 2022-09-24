package main

import (
	"github.com/palomachain/pigeon/app"
	"github.com/spf13/cobra"
)

var (
	healthCheckVerify = &cobra.Command{
		Use:   "health-check",
		Short: "Verifies the health of pigeon.",
		RunE: func(cmd *cobra.Command, args []string) error {
			app.HealthCheckService().BootChecker(cmd.Context())
			app.HealthCheckService().Check(cmd.Context())
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(healthCheckVerify)
	configRequired(healthCheckVerify)
}
