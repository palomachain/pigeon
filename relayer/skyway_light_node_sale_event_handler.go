package relayer

import (
	"context"
	"sync"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/internal/liblog"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) SkywayHandleLightNodeSaleEvent(ctx context.Context, locker sync.Locker) error {
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
	err = r.handleLightNodeSaleEvents(ctx, r.processors)
	locker.Unlock()

	return handleProcessError(ctx, err)
}

func (r *Relayer) handleLightNodeSaleEvents(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	for _, p := range processors {
		chainReferenceID := p.GetChainReferenceID()

		logger := liblog.WithContext(ctx).WithFields(log.Fields{
			"chain-reference-id": chainReferenceID,
			"action":             "handle-skyway-light-node-sale-events",
		})

		batchSaleEvents, err := p.GetLightNodeSaleEvents(ctx, r.palomaClient.GetCreator())
		if err != nil {
			logger.WithError(err).Error("couldn't get events")
			return err
		}

		logger = logger.WithFields(log.Fields{
			"event-nonces": slice.Map(batchSaleEvents, func(event chain.LightNodeSaleEvent) uint64 {
				return event.EventNonce
			}),
			"skyway-nonces": slice.Map(batchSaleEvents, func(event chain.LightNodeSaleEvent) uint64 {
				return event.SkywayNonce
			}),
		})

		logger.Debug("got ", len(batchSaleEvents), " events")
		if len(batchSaleEvents) > 0 {
			logger.Info("claiming for ", len(batchSaleEvents), " events")
			err := p.SubmitLightNodeSaleClaims(ctx, batchSaleEvents, r.palomaClient.GetCreator())
			if err != nil {
				logger.WithError(err).Error("error submitting claim for events")
				return err
			}
		}
	}

	return nil
}
