package relayer

import (
	"context"
	"sync"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) SignMessages(ctx context.Context, _ sync.Locker) error {
	log.Info("signer loop")
	if ctx.Err() != nil {
		log.Info("exiting signer loop as context has ended")
		return ctx.Err()
	}

	err := r.buildProcessors(ctx, nil)
	if err != nil {
		log.Error(err)
		return err
	}

	err = r.signMessages(ctx, r.processors)

	return handleProcessError(ctx, err)
}

func (r *Relayer) signMessages(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	// todo randomise
	for _, p := range processors {
		// todo randomise
		for _, queueName := range p.SupportedQueues() {
			logger := log.WithFields(log.Fields{
				"queue-name": queueName,
				"action":     "sign",
			})

			messagesForSigning, err := r.palomaClient.QueryMessagesForSigning(ctx, queueName)
			if err != nil {
				logger.Error("failed getting messages to sign")
				if isFatal(err) {
					return err
				}
				// Move on to the next queue on the same chain
				continue
			}

			logger = logger.WithFields(log.Fields{
				"message-ids": slice.Map(messagesForSigning, func(msg chain.QueuedMessage) uint64 {
					return msg.ID
				}),
			})

			if len(messagesForSigning) > 0 {
				logger := logger.WithFields(log.Fields{
					"messages-to-sign": slice.Map(messagesForSigning, func(msg chain.QueuedMessage) uint64 {
						return msg.ID
					}),
				})

				logger.Info("signing ", len(messagesForSigning), " messages")
				signedMessages, err := p.SignMessages(ctx, messagesForSigning...)
				if err != nil {
					logger.WithError(err).Error("unable to sign messages")
					// If we fail to sign this batch, we will fail to sign them
					// all, so might as well return now
					return err
				}
				logger = logger.WithFields(log.Fields{
					"signed-messages": slice.Map(signedMessages, func(msg chain.SignedQueuedMessage) log.Fields {
						return log.Fields{
							"id": msg.ID,
						}
					}),
				})
				logger.Info("signed messages")

				if err = r.broadcastSignatures(ctx, queueName, signedMessages); err != nil {
					logger.WithError(err).Error("couldn't broadcast signatures and process attestation")
					if isFatal(err) {
						return err
					}
				}
			}

		}
	}

	return nil
}

func (r *Relayer) broadcastSignatures(ctx context.Context, queueTypeName string, sigs []chain.SignedQueuedMessage) error {
	broadcastMessageSignatures, err := slice.MapErr(
		sigs,
		func(sig chain.SignedQueuedMessage) (paloma.BroadcastMessageSignatureIn, error) {
			log.WithFields(
				log.Fields{
					"id":              sig.ID,
					"queue-type-name": queueTypeName,
				},
			).Debug("broadcasting signed message")

			return paloma.BroadcastMessageSignatureIn{
				ID:              sig.ID,
				QueueTypeName:   queueTypeName,
				Signature:       sig.Signature,
				SignedByAddress: sig.SignedByAddress,
			}, nil
		},
	)
	if err != nil {
		return err
	}

	return r.palomaClient.BroadcastMessageSignatures(ctx, broadcastMessageSignatures...)
}
