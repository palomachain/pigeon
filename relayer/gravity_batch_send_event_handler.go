package relayer

import (
	"context"
	"sync"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/internal/liblog"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) GravityHandleBatchSendEvent(ctx context.Context, locker sync.Locker) error {
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
	err = r.handleBatchSendEvents(ctx, r.processors)
	locker.Unlock()

	return handleProcessError(ctx, err)
}

func (r *Relayer) handleBatchSendEvents(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	for _, p := range processors {
		chainReferenceID := p.GetChainReferenceID()

		logger := liblog.WithContext(ctx).WithFields(log.Fields{
			"chain-reference-id": chainReferenceID,
			"action":             "handle-gravity-batch-send-events",
		})

		batchSendEvents, err := p.GetBatchSendEvents(ctx, r.palomaClient.GetCreator())
		if err != nil {
			logger.WithError(err).Error("couldn't get events")
			return err
		}

		logger = logger.WithFields(log.Fields{
			"event-nonces": slice.Map(batchSendEvents, func(event chain.BatchSendEvent) uint64 {
				return event.EventNonce
			}),
			"gravity-nonces": slice.Map(batchSendEvents, func(event chain.BatchSendEvent) uint64 {
				return event.GravityNonce
			}),
		})

		logger.Debug("got ", len(batchSendEvents), " events")
		if len(batchSendEvents) > 0 {
			// Walk through the different batchSendEvents and do different things for different batchSendEvents

			logger.Info("claiming for ", len(batchSendEvents), " events")
			err := p.SubmitBatchSendToEthClaims(ctx, batchSendEvents, r.palomaClient.GetCreator())
			if err != nil {
				logger.WithError(err).Error("error submitting claim for events")
				return err
			}
		}
	}

	return nil
}
