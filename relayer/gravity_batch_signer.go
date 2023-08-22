package relayer

import (
	"context"
	gravity "github.com/palomachain/paloma/x/gravity/types"
	"sync"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) GravitySignBatches(ctx context.Context, locker sync.Locker) error {
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
	err = r.gravitySignBatches(ctx, r.processors)
	locker.Unlock()

	return handleProcessError(err)
}

func (r *Relayer) gravitySignBatches(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	for _, p := range processors {
		chainReferenceID := p.GetChainReferenceID()

		logger := log.WithFields(log.Fields{
			"chain-reference-id": chainReferenceID,
			"action":             "sign-gravity-batches",
		})

		// Get all batches that we haven't signed
		batchesForSigning, err := r.palomaClient.GravityQueryLastUnsignedBatch(ctx, chainReferenceID)
		if err != nil {
			logger.WithError(err).Error("failed getting batches to sign")
			return err
		}

		logger = logger.WithFields(log.Fields{
			"batch-nonces": slice.Map(batchesForSigning, func(batch gravity.OutgoingTxBatch) uint64 {
				return batch.BatchNonce
			}),
		})

		if len(batchesForSigning) > 0 {
			logger.Info("signing ", len(batchesForSigning), " batches")
			signedBatches, err := p.GravitySignBatches(ctx, batchesForSigning...)
			if err != nil {
				logger.WithError(err).Error("unable to sign batches")
				return err
			}
			logger = logger.WithFields(log.Fields{
				"signed-batches": slice.Map(signedBatches, func(batch chain.SignedGravityOutgoingTxBatch) log.Fields {
					return log.Fields{
						"nonce": batch.BatchNonce,
					}
				}),
			})
			logger.Info("signed batches")

			if err = r.gravityConfirmBatches(ctx, signedBatches); err != nil {
				logger.WithError(err).Error("couldn't broadcast signatures and process attestation")
				return err
			}

		}
	}

	return nil
}

func (r *Relayer) gravityConfirmBatches(ctx context.Context, sigs []chain.SignedGravityOutgoingTxBatch) error {
	return r.palomaClient.GravityConfirmBatches(ctx, sigs...)
}
