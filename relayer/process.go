package relayer

import (
	"context"
	"fmt"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/chain/paloma/collision"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) Process(ctx context.Context, processors []chain.Processor) error {
	ctx, cleanup := r.zeroCollisionStrategy.GoStartLane(ctx)
	defer cleanup()

	// todo randomise
	for _, p := range processors {

		// todo randomise
		for _, queueName := range p.SupportedQueues() {
			logger := log.WithFields(log.Fields{
				"queue-name": queueName,
			})

			// TODO: remove comments once signing is done on the paloma side.
			queuedMessages, err := r.palomaClient.QueryMessagesForSigning(ctx, queueName)
			loggerQueuedMessages := logger.WithFields(log.Fields{
				"message-ids": slice.Map(queuedMessages, func(msg chain.QueuedMessage) uint64 {
					return msg.ID
				}),
			})

			if err != nil {
				logger.Warn("failed getting messages to sign")
				return err
			}

			if len(queuedMessages) > 0 {
				loggerQueuedMessages.Info("messages to sign")
				signedMessages, err := p.SignMessages(ctx, queueName, queuedMessages...)
				if err != nil {
					loggerQueuedMessages.WithFields(log.Fields{
						"err": err,
					}).Error("unable to sign messages")
					return err
				}
				loggerQueuedMessages = loggerQueuedMessages.WithFields(log.Fields{
					"signed-messages": slice.Map(signedMessages, func(msg chain.SignedQueuedMessage) log.Fields {
						return log.Fields{
							"id":        msg.ID,
							"signature": msg.Signature,
						}
					}),
				})
				loggerQueuedMessages.Info("signed messages")

				if err = r.broadcastSignatures(ctx, queueName, signedMessages); err != nil {
					loggerQueuedMessages.WithFields(log.Fields{
						"err": err,
					}).Info("couldn't broadcast signatures and process attestation")
					return err
				}
			}

			relayCandidateMsgs, err := r.palomaClient.QueryMessagesInQueue(ctx, queueName)

			logger.Debug("got 246")

			relayCandidateMsgs = slice.Filter(relayCandidateMsgs, func(msg chain.MessageWithSignatures) bool {
				return collision.AllowedToExecute(
					ctx,
					fmt.Sprintf("%s-%d", queueName, msg.ID),
					r.palomaClient.GetValidatorAddress(),
				)
			})

			if err != nil {
				logger.WithFields(log.Fields{
					"err": err,
				}).Error("couldn't get messages to relay")
				return err
			}

			logger = logger.WithFields(log.Fields{
				"messages-to-relay": slice.Map(relayCandidateMsgs, func(msg chain.MessageWithSignatures) uint64 {
					return msg.ID
				}),
			})

			if len(relayCandidateMsgs) > 0 {
				logger.Info("relaying messages")
				if err = p.ProcessMessages(ctx, queueName, relayCandidateMsgs); err != nil {
					logger.WithField("err", err).Error("error relaying messages")
					return err
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
