package relayer

import (
	"context"
	"sync"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) SkywayEstimateBatchGas(ctx context.Context, locker sync.Locker) error {
	log.Info("gas estimation loop")
	if ctx.Err() != nil {
		log.Info("exiting gas estimation loop as context has ended")
		return ctx.Err()
	}

	err := r.buildProcessors(ctx, locker)
	if err != nil {
		log.Error(err)
		return err
	}

	locker.Lock()
	err = r.skywayEstimateBatchGas(ctx, r.processors)
	locker.Unlock()

	return handleProcessError(ctx, err)
}

func (r *Relayer) skywayEstimateBatchGas(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	for _, p := range processors {
		chainReferenceID := p.GetChainReferenceID()

		logger := log.WithFields(log.Fields{
			"chain-reference-id": chainReferenceID,
			"action":             "gas-estimate-skyway-batches",
		})

		// Get all batches that we haven't estimated
		batchesForEstimating, err := r.palomaClient.SkywayQueryLastPendingBatchForGasEstimation(ctx, chainReferenceID)
		if err != nil {
			logger.WithError(err).Error("failed getting batches to estimate")
			return err
		}

		logger = logger.WithFields(log.Fields{
			"batch-nonces": slice.Map(batchesForEstimating, func(batch chain.SkywayBatchWithSignatures) uint64 {
				return batch.BatchNonce
			}),
		})

		if len(batchesForEstimating) > 0 {
			logger.Info("estimating ", len(batchesForEstimating), " batches")
			estimatedBatches, err := p.SkywayEstimateBatches(ctx, batchesForEstimating)
			if err != nil {
				logger.WithError(err).Error("unable to estimate batches")
				return err
			}
			logger = logger.WithFields(log.Fields{
				"estimated-batches": slice.Map(estimatedBatches, func(batch chain.EstimatedSkywayBatch) log.Fields {
					return log.Fields{
						"nonce": batch.BatchNonce,
					}
				}),
			})
			logger.Info("estimated batches")

			if err = r.palomaClient.SkywayEstimateBatchGas(ctx, estimatedBatches...); err != nil {
				logger.WithError(err).Error("couldn't broadcast gas estimates for batch.")
				return err
			}

		}
	}

	return nil
}
