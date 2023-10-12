package relayer

import (
	"context"
	"sync"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/internal/liblog"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) GravityRelayBatches(ctx context.Context, locker sync.Locker) error {
	logger := liblog.WithContext(ctx)
	logger.Info("relayer loop")
	if ctx.Err() != nil {
		logger.Info("exiting relayer loop as context has ended")
		return ctx.Err()
	}

	err := r.buildProcessors(ctx, locker)
	if err != nil {
		return err
	}

	locker.Lock()
	err = r.gravityRelayBatches(ctx, r.processors)
	locker.Unlock()

	return handleProcessError(ctx, err)
}

func (r *Relayer) gravityRelayBatches(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	for _, p := range processors {
		chainReferenceID := p.GetChainReferenceID()

		logger := liblog.WithContext(ctx).WithFields(log.Fields{
			"chain-reference-id": chainReferenceID,
			"action":             "relay-gravity-batches",
		})

		batchesForRelaying, err := r.palomaClient.GravityQueryBatchesForRelaying(ctx, chainReferenceID)

		logger = logger.WithFields(log.Fields{
			"batch-nonces": slice.Map(batchesForRelaying, func(batch chain.GravityBatchWithSignatures) uint64 {
				return batch.BatchNonce
			}),
		})

		logger.Debug("got ", len(batchesForRelaying), " batches")
		if err != nil {
			logger.WithError(err).Error("couldn't get batches to relay")
			return err
		}

		if len(batchesForRelaying) > 0 {
			logger.Info("relaying ", len(batchesForRelaying), " batches")
			err := p.GravityRelayBatches(ctx, batchesForRelaying)
			if err != nil {
				logger.WithError(err).Error("error relaying batches")
				return err
			}
		}
	}
	return nil
}
