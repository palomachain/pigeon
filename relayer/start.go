package relayer

import (
	"context"
	"time"

	"github.com/palomachain/pigeon/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vizualni/whoops"
)

const (
	defaultErrorCountToExit = 5

	defaultLoopTimeout = 1 * time.Minute
)

// Start starts the relayer. It's responsible for handling the communication
// with Paloma and other chains.
func (r *Relayer) Start(ctx context.Context) error {
	log.Info("starting relayer")

	if err := r.init(); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("couldn't initialize relayer")
		return err
	}
	processors, err := r.buildProcessors(ctx)
	if err != nil {
		return err
	}

	if err := r.updateExternalChainInfos(ctx, processors); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("couldn't update external chain info")
		return err
	}

	consecutiveFailures := whoops.Group{}

	timer := time.NewTicker(defaultLoopTimeout)

	// only used to enter into the loop below emmidiaetly after the first "tick"
	firstLoopEnter := make(chan struct{}, 1)
	firstLoopEnter <- struct{}{}

	process := func() error {
		log.Info("relayer loop")

		processors, err := r.buildProcessors(ctx)
		if err != nil {
			return err
		}

		err = r.Process(ctx, processors)
		if err == nil {
			// resetting the failures
			if len(consecutiveFailures) > 0 {
				consecutiveFailures = whoops.Group{}
			}
			return nil
		}

		if errors.IsUnrecoverable(err) {
			// there is no way that we can recover from this
			log.WithFields(log.Fields{
				"err": err,
			}).Error("unrecoverable error returned")
			return err
		}

		consecutiveFailures.Add(err)

		log.WithFields(log.Fields{
			"err": err,
		}).Warn("adding error to consecutive failures")

		if len(consecutiveFailures) >= defaultErrorCountToExit {
			log.Error("too many consecutive failures")
			return errors.Unrecoverable(consecutiveFailures)
		}

		log.WithFields(log.Fields{
			"err": err,
		}).Error("error while trying to relay messages")

		return nil
	}

	for {
		log.Debug("waiting on the loop for a new tick")
		select {
		case <-ctx.Done():
			log.Warn("exiting due to context being done")
			return ctx.Err()
		case <-firstLoopEnter:
			// don't put anything into the firstLoopEnter channel anymore
			if err := process(); err != nil {
				return err
			}
		case <-timer.C:
			if err := process(); err != nil {
				return err
			}
		}
	}
}
