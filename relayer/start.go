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

	defaultLoopTimeout = 10 * time.Second
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

	for {
		select {
		case <-ctx.Done():
			log.Warn("exiting due to context being done")
			return ctx.Err()
		case <-time.After(defaultLoopTimeout):
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
				continue
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
				log.WithFields(log.Fields{
					"err": consecutiveFailures,
				}).Error("too many consecutive failures")
				return errors.Unrecoverable(consecutiveFailures)
			}

			log.WithFields(log.Fields{
				"err": err,
			}).Error("error while trying to relay messages")
		}
	}
}
