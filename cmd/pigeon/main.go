package main

import (
	"os"
	"strings"

	"github.com/palomachain/pigeon/app"
	"github.com/sirupsen/logrus"
)

const (
	logLevelEnvName  = app.AppNameCaps + "_LOG_LEVEL"
	logFormatEnvName = app.AppNameCaps + "_LOG_FORMAT"
)

func main() {
	level, err := logrus.ParseLevel(strings.TrimSpace(os.Getenv(logLevelEnvName)))
	if err == nil {
		logrus.SetLevel(level)
	}

	if level == logrus.TraceLevel {
		logrus.SetReportCaller(true)
	}

	formatter := strings.ToLower(strings.TrimSpace(os.Getenv(logFormatEnvName)))

	if formatter == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
