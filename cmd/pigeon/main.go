package main

import (
	"os"
	"strings"

	"github.com/palomachain/pigeon/app"
	"github.com/sirupsen/logrus"
)

const (
	logLevelEnvName = app.AppNameCaps + "_LOG_LEVEL"
)

func main() {
	level, err := logrus.ParseLevel(strings.TrimSpace(os.Getenv(logLevelEnvName)))
	if err == nil {
		logrus.SetLevel(level)
	}

	if level == logrus.TraceLevel {
		logrus.SetReportCaller(true)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
