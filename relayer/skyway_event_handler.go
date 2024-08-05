package relayer

import (
	"context"
	"sync"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/internal/liblog"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) SkywayHandleEvents(ctx context.Context, locker sync.Locker) error {
	logger := liblog.WithContext(ctx)
	logger.Info("event watcher loop")
	if ctx.Err() != nil {
		logger.Info("exiting event watcher loop as context has ended")
		return ctx.Err()
	}

	err := r.buildProcessors(ctx, locker)
	if err != nil {
		return err
	}

	locker.Lock()
	err = r.handleEvents(ctx, r.processors)
	locker.Unlock()

	return handleProcessError(ctx, err)
}

func (r *Relayer) handleEvents(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	for _, p := range processors {
		chainReferenceID := p.GetChainReferenceID()

		logger := liblog.WithContext(ctx).WithFields(log.Fields{
			"chain-reference-id": chainReferenceID,
			"action":             "handle-skyway-events",
		})

		events, err := p.GetSkywayEvents(ctx, r.palomaClient.GetCreator())
		if err != nil {
			logger.WithError(err).Error("couldn't get events")
			return err
		}

		logger = logger.WithFields(log.Fields{
			"event-nonces": slice.Map(events, func(event chain.SkywayEventer) uint64 {
				return event.GetEventNonce()
			}),
			"skyway-nonces": slice.Map(events, func(event chain.SkywayEventer) uint64 {
				return event.GetSkywayNonce()
			}),
		})

		logger.Debug("got ", len(events), " events")

		if len(events) > 0 {
			logger.Info("claiming for ", len(events), " events")

			err := p.SubmitEventClaims(ctx, events, r.palomaClient.GetCreator())
			if err != nil {
				logger.WithError(err).Error("error submitting claim for events")
				return err
			}
		}
	}

	return nil
}
