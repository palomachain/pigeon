package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/palomachain/pigeon/app"
	"github.com/spf13/cobra"
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "starts the pigeon server",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("*chirp chirp*")
			ctx := catchKillSignal(cmd.Context(), 30*time.Second)
			err := app.Relayer().Start(ctx)
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return err
		},
	}
)

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
