package main

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	logLevelEnvName = "LOG_LEVEL"
)

func main() {
	level, err := logrus.ParseLevel(strings.TrimSpace(os.Getenv(logLevelEnvName)))
	if err == nil {
		logrus.SetLevel(level)
	}

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
