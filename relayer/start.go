package relayer

import (
	"context"
	goerrors "errors"
	"github.com/palomachain/pigeon/errors"
	"sync"
	"time"

	"github.com/VolumeFi/whoops"
	"github.com/palomachain/pigeon/util/channels"
	log "github.com/sirupsen/logrus"
)

const (
	defaultErrorCountToExit = 5

	defaultLoopTimeout = 1 * time.Minute
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

// Start starts the relayer. It's responsible for handling the communication
// with Paloma and other chains.
func (r *Relayer) Start(ctx context.Context) error {
	if err := r.waitUntilStaking(ctx); err != nil {
		return err
	}

	log.Info("starting relayer")
	var locker sync.Mutex

	go r.startUpdateExternalChainInfos(ctx, &locker)

	ticker := time.NewTicker(defaultLoopTimeout)
	defer ticker.Stop()

	// only used to enter into the loop below immediately after the first "tick"
	firstLoopEnter := make(chan time.Time, 1)
	firstLoopEnter <- time.Time{}

	go func() {
		r.startKeepAlive(ctx, &locker)
	}()

	tickerCh := channels.FanIn(ticker.C, firstLoopEnter)
	for {
		log.Debug("waiting on the loop for a new tick")
		select {
		case <-ctx.Done():
			log.Warn("exiting due to context being done")
			return ctx.Err()
		case _, chOpen := <-tickerCh:
			if !chOpen {
				if ctx.Err() != nil {
					return nil
				}
				return whoops.WrapS(ErrUnknown, "ticker channel for message processing was closed unexpectedly")
			}
			if err := r.process(ctx, &locker); err != nil {
				log.WithError(err).Error("error while trying to process messages")
			}
		}
	}
}

func (r *Relayer) startUpdateExternalChainInfos(ctx context.Context, locker sync.Locker) {
	ticker := time.NewTicker(defaultLoopTimeout)
	defer ticker.Stop()

	for range ticker.C {
		processors, err := r.buildProcessors(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("couldn't build processors to update external chain info")

			continue
		}

		log.Info("trying to update external chain info")

		locker.Lock()
		err = r.updateExternalChainInfos(ctx, processors)
		locker.Unlock()

		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("couldn't update external chain info. Will try again.")
		}
	}
}

func (r *Relayer) process(ctx context.Context, locker sync.Locker) error {
	log.Info("relayer loop")
	if ctx.Err() != nil {
		log.Info("exiting relayer loop as context has ended")
		return ctx.Err()
	}

	processors, err := r.buildProcessors(ctx)
	if err != nil {
		return err
	}

	locker.Lock()
	err = r.Process(ctx, processors)
	locker.Unlock()

	switch {
	case err == nil:
		// success
		return nil
	case goerrors.Is(err, context.Canceled):
		log.WithFields(log.Fields{
			"err": err,
		}).Debug("exited from the process loop due the context being canceled")
		return nil
	case goerrors.Is(err, context.DeadlineExceeded):
		log.WithFields(log.Fields{
			"err": err,
		}).Debug("exited from the process loop due the context deadline being exceeded")
		return nil
	case errors.IsUnrecoverable(err):
		// there is no way that we can recover from this
		log.WithFields(log.Fields{
			"err": err,
		}).Error("unrecoverable error returned")
		return err
	default:
		log.WithFields(log.Fields{
			"err": err,
		}).Error("error returned in process loop")
		return nil
	}
}
