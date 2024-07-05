package relayer

import (
	"context"
	"sync"

	skyway "github.com/palomachain/paloma/x/skyway/types"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) SkywaySignBatches(ctx context.Context, locker sync.Locker) error {
	log.Info("signer loop")
	if ctx.Err() != nil {
		log.Info("exiting signer loop as context has ended")
		return ctx.Err()
	}

	err := r.buildProcessors(ctx, locker)
	if err != nil {
		log.Error(err)
		return err
	}

	locker.Lock()
	err = r.skywaySignBatches(ctx, r.processors)
	locker.Unlock()

	return handleProcessError(ctx, err)
}

func (r *Relayer) skywaySignBatches(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	for _, p := range processors {
		chainReferenceID := p.GetChainReferenceID()

		logger := log.WithFields(log.Fields{
			"chain-reference-id": chainReferenceID,
			"action":             "sign-skyway-batches",
		})

		// Get all batches that we haven't signed
		batchesForSigning, err := r.palomaClient.SkywayQueryLastUnsignedBatch(ctx, chainReferenceID)
		if err != nil {
			logger.WithError(err).Error("failed getting batches to sign")
			return err
		}

		logger = logger.WithFields(log.Fields{
			"batch-nonces": slice.Map(batchesForSigning, func(batch skyway.OutgoingTxBatch) uint64 {
				return batch.BatchNonce
			}),
		})

		if len(batchesForSigning) > 0 {
			logger.Info("signing ", len(batchesForSigning), " batches")
			signedBatches, err := p.SkywaySignBatches(ctx, batchesForSigning...)
			if err != nil {
				logger.WithError(err).Error("unable to sign batches")
				return err
			}
			logger = logger.WithFields(log.Fields{
				"signed-batches": slice.Map(signedBatches, func(batch chain.SignedSkywayOutgoingTxBatch) log.Fields {
					return log.Fields{
						"nonce": batch.BatchNonce,
					}
				}),
			})
			logger.Info("signed batches")

			if err = r.skywayConfirmBatches(ctx, signedBatches); err != nil {
				logger.WithError(err).Error("couldn't broadcast signatures and process attestation")
				return err
			}

		}
	}

	return nil
}

func (r *Relayer) skywayConfirmBatches(ctx context.Context, sigs []chain.SignedSkywayOutgoingTxBatch) error {
	return r.palomaClient.SkywayConfirmBatches(ctx, sigs...)
}
