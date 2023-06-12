package relayer

import (
	"context"
	goerrors "errors"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	updateExternalChainsLoopInterval = 1 * time.Minute
	signMessagesLoopInterval         = 1 * time.Minute
	relayMessagesLoopInterval        = 1 * time.Minute
	attestMessagesLoopInterval       = 1 * time.Minute
)

func (r *Relayer) waitUntilStaking(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err := r.isStaking(ctx)
		switch {
		case err == nil:
			// validator is staking
			log.Info("validator is staking")
			return nil
		case goerrors.Is(err, ErrValidatorIsNotStaking):
			// does nothing
			log.Info("not staking. waiting")
		default:
			return err
		}

		time.Sleep(5 * time.Second)
	}
}

func (r *Relayer) startProcess(ctx context.Context, locker sync.Locker, tickerInterval time.Duration, process func(context.Context, sync.Locker) error) {
	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	logger := log.WithFields(log.Fields{})
	for {
		select {
		case <-ctx.Done():
			logger.Warn("exiting due to context being done")
			return
		case <-ticker.C:
			err := process(ctx, locker)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

// Start starts the relayer. It's responsible for handling the communication
// with Paloma and other chains.
func (r *Relayer) Start(ctx context.Context) error {
	if err := r.waitUntilStaking(ctx); err != nil {
		return err
	}

	log.Info("starting relayer")
	var locker sync.Mutex

	// Start background goroutines to run separately from each other
	go r.startProcess(ctx, &locker, updateExternalChainsLoopInterval, r.UpdateExternalChainInfos)
	go r.startProcess(ctx, &locker, signMessagesLoopInterval, r.SignMessages)
	go r.startProcess(ctx, &locker, relayMessagesLoopInterval, r.RelayMessages)
	go r.startProcess(ctx, &locker, attestMessagesLoopInterval, r.AttestMessages)

	// Start the foreground process
	r.startProcess(ctx, &locker, r.relayerConfig.KeepAliveLoopTimeout, r.keepAlive)
	return nil
}
