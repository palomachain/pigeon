package relayer

import (
	"context"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) Process(ctx context.Context, processors []chain.Processor) error {
	for chainID, p := range processors {
		for _, queueName := range p.SupportedQueues() {
			logger := log.WithFields(log.Fields{
				"processor-chain-id": chainID,
				"queue-name":         queueName,
			})

			// TODO: remove comments once signing is done on the paloma side.
			queuedMessages, err := r.palomaClient.QueryMessagesForSigning(ctx, queueName)
			loggerQueuedMessages := logger.WithFields(log.Fields{
				"message-ids": slice.Map(queuedMessages, func(msg chain.QueuedMessage) uint64 {
					return msg.ID
				}),
			})
			loggerQueuedMessages.Info("messages to sign")

			if err != nil {
				logger.Warn("failed getting messages to sign")
				return err
			}

			if len(queuedMessages) > 0 {
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

				if err = r.broadcastSignaturesAndProcessAttestation(ctx, queueName, signedMessages); err != nil {
					loggerQueuedMessages.WithFields(log.Fields{
						"err": err,
					}).Info("couldn't broadcast signatures and process attestation")
					return err
				}
			}

			relayCandidateMsgs, err := r.palomaClient.QueryMessagesInQueue(ctx, queueName)
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

			logger.Info("relaying messages")
			if err = p.ProcessMessages(ctx, queueName, relayCandidateMsgs); err != nil {
				logger.WithField("err", err).Error("error relaying messages")
				return err
			}
		}
	}

	return nil
}

func (r *Relayer) broadcastSignaturesAndProcessAttestation(ctx context.Context, queueTypeName string, sigs []chain.SignedQueuedMessage) error {
	broadcastMessageSignatures, err := slice.MapErr(
		sigs,
		func(sig chain.SignedQueuedMessage) (paloma.BroadcastMessageSignatureIn, error) {
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
