package relayer

import (
	"context"
	"errors"
	"fmt"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/chain/paloma/collision"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) Process(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	ctx, cleanup, err := collision.GoStartLane(ctx, r.palomaClient, r.palomaClient.GetValidatorAddress())
	if err != nil {
		return err
	}
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
				if errors.Is(err, context.DeadlineExceeded) {
					return nil
				}
				logger.Error("failed getting messages to sign")
				return err
			}

			if len(queuedMessages) > 0 {
				loggerQueuedMessages.Info("messages to sign")
				signedMessages, err := p.SignMessages(ctx, queueName, queuedMessages...)
				if errors.Is(err, context.DeadlineExceeded) {
					return nil
				}
				if err != nil {
					loggerQueuedMessages.WithError(err).Error("unable to sign messages")
					return err
				}
				loggerQueuedMessages = loggerQueuedMessages.WithFields(log.Fields{
					"signed-messages": slice.Map(signedMessages, func(msg chain.SignedQueuedMessage) log.Fields {
						return log.Fields{
							"id": msg.ID,
						}
					}),
				})
				loggerQueuedMessages.Info("signed messages")

				if err = r.broadcastSignatures(ctx, queueName, signedMessages); err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						return nil
					}
					loggerQueuedMessages.WithError(err).Error("couldn't broadcast signatures and process attestation")
					return err
				}
			}

			msgsInQueue, err := r.palomaClient.QueryMessagesInQueue(ctx, queueName)

			logger.Debug("got ", len(msgsInQueue), " messages from ", queueName)

			relayCandidateMsgs := slice.Filter(
				msgsInQueue,
				func(msg chain.MessageWithSignatures) bool {
					return len(msg.PublicAccessData) == 0
				},
				func(msg chain.MessageWithSignatures) bool {
					return collision.AllowedToExecute(
						ctx,
						[]byte(fmt.Sprintf("%s-%d", queueName, msg.ID)),
					)
				},
			)

			msgsToProvideEvidenceFor := slice.Filter(msgsInQueue, func(msg chain.MessageWithSignatures) bool {
				return len(msg.PublicAccessData) > 0
			})

			if err != nil {
				logger.WithError(err).Error("couldn't get messages to relay")
				return err
			}

			if len(relayCandidateMsgs) > 0 {
				logger := logger.WithFields(log.Fields{
					"messages-to-relay": slice.Map(relayCandidateMsgs, func(msg chain.MessageWithSignatures) uint64 {
						return msg.ID
					}),
				})
				logger.Info("relaying messages")
				if err = p.ProcessMessages(ctx, queueName, relayCandidateMsgs); err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						return nil
					}
					logger.WithFields(log.Fields{
						"err":        err,
						"queue-name": queueName,
						"messages-to-relay": slice.Map(relayCandidateMsgs, func(msg chain.MessageWithSignatures) uint64 {
							return msg.ID
						}),
					}).Error("error relaying messages")
					return err
				}
			}

			if len(msgsToProvideEvidenceFor) > 0 {
				logger := logger.WithFields(log.Fields{
					"messages-to-provide-evidence-for": slice.Map(msgsToProvideEvidenceFor, func(msg chain.MessageWithSignatures) uint64 {
						return msg.ID
					}),
				})
				logger.Info("providing evidence for messages")
				if err = p.ProvideEvidence(ctx, queueName, msgsToProvideEvidenceFor); err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						return nil
					}
					logger.WithError(err).Error("error providing evidence for messages")
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
