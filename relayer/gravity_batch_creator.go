package relayer

import (
	"context"
	"sync"

	"github.com/palomachain/pigeon/chain"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) GravityCreateBatches(ctx context.Context, locker sync.Locker) error {
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
	err = r.gravityCreateBatches(ctx, r.processors)
	locker.Unlock()

	return handleProcessError(err)
}

func (r *Relayer) gravityCreateBatches(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	for _, p := range processors {
		chainReferenceID := p.GetChainReferenceID()

		logger := log.WithFields(log.Fields{
			"chain-reference-id": chainReferenceID,
			"action":             "create-gravity-batches",
		})

		// First, request that Gravity create a batch
		err := r.palomaClient.GravityRequestBatch(ctx, chainReferenceID)
		if err != nil {
			logger.Error("failed requesting batch")
			return err
		}
	}
	return nil
}
