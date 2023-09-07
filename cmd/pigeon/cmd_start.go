package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/palomachain/pigeon/app"
	"github.com/palomachain/pigeon/health"
	"github.com/palomachain/pigeon/internal/mev"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts the pigeon server",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := os.Getpid()
		fmt.Println("*coo! coo!*")
		log.WithFields(
			log.Fields{
				"PID":     pid,
				"version": app.Version(),
				"commit":  app.Commit(),
			},
		).Info("app info")

		ctx := catchKillSignal(cmd.Context(), 30*time.Second)

		// start healthcheck server
		go func() {
			health.StartHTTPServer(
				ctx,
				app.Config().HealthCheckAddress,
				app.Config().HealthCheckPort,
				pid,
				app.Version(),
				app.Commit(),
			)
		}()

		// wait for paloma to get online
		waitCtx, cancelFnc := context.WithTimeout(ctx, 2*time.Minute)
		err := health.WaitForPaloma(waitCtx, app.PalomaClient())
		cancelFnc()
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			log.WithError(err).Fatal("exiting as paloma was not detected to be running")
			return err
		}

		// build a context that will get canceled if paloma ever goes offline
		ctx = health.CancelContextIfPalomaIsDown(ctx, app.PalomaClient())

		relayer := app.Relayer()
		relayer.SetAppVersion(app.Version())
		relayer.SetMevClient(mev.New(app.Config()))

		err = relayer.Start(ctx)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	configRequired(startCmd)
}

func catchKillSignal(ctx context.Context, waitTimeout time.Duration) context.Context {
	retCtx, closeCtx := context.WithCancel(ctx)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, os.Kill)
	nextSignalShouldKillTheProcess := false
	go func() {
		for range signalCh {
			if nextSignalShouldKillTheProcess {
				fmt.Println("exiting forcefully")
				os.Exit(1)
			}
			fmt.Printf("pigeon will close in %s max\n", waitTimeout)
			fmt.Println("press ctrl+c again to forcefully close pigeon (not recommended)")
			closeCtx()
			nextSignalShouldKillTheProcess = true
			go func() {
				<-time.NewTimer(waitTimeout).C
				fmt.Printf("pigeon didn't properly close in %s. forcing exit", waitTimeout)
				os.Exit(1)
			}()
		}
	}()

	return retCtx
}
