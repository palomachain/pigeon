package relayer

import (
	"context"
	"fmt"

	"github.com/VolumeFi/whoops"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/chain/paloma/collision"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) SignMessages(ctx context.Context, queueName string, messagesForSigning []chain.QueuedMessage, processor chain.Processor) error {
	loggerQueuedMessages := log.WithFields(log.Fields{
		"queue-name": queueName,
		"message-ids": slice.Map(messagesForSigning, func(msg chain.QueuedMessage) uint64 {
			return msg.ID
		}),
	})

	if len(messagesForSigning) > 0 {
		loggerQueuedMessages.Info("messages to sign")
		signedMessages, err := processor.SignMessages(ctx, queueName, messagesForSigning...)
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
			loggerQueuedMessages.WithError(err).Error("couldn't broadcast signatures and process attestation")
			return err
		}
	}
	return nil
}

func (r *Relayer) RelayMessages(ctx context.Context, queueName string, messagesInQueue []chain.MessageWithSignatures, p chain.Processor) error {
	logger := log.WithFields(log.Fields{
		"queue-name": queueName,
	})

	ctx, cleanup, err := collision.GoStartLane(ctx, r.palomaClient, r.palomaClient.GetValidatorAddress())
	if err != nil {
		return err
	}
	defer cleanup()

	relayCandidateMsgs := slice.Filter(
		messagesInQueue,
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

	if len(relayCandidateMsgs) > 0 {
		logger := logger.WithFields(log.Fields{
			"messages-to-relay": slice.Map(relayCandidateMsgs, func(msg chain.MessageWithSignatures) uint64 {
				return msg.ID
			}),
		})
		logger.Info("relaying messages")
		err := p.ProcessMessages(ctx, queueName, relayCandidateMsgs)
		if err != nil {
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
	return nil
}

func (r *Relayer) ProvideEvidenceForMessages(ctx context.Context, queueName string, messagesInQueue []chain.MessageWithSignatures, p chain.Processor) error {
	logger := log.WithFields(log.Fields{
		"queue-name": queueName,
	})

	msgsToProvideEvidenceFor := slice.Filter(messagesInQueue, func(msg chain.MessageWithSignatures) bool {
		return len(msg.PublicAccessData) > 0
	})

	if len(msgsToProvideEvidenceFor) > 0 {
		logger := logger.WithFields(log.Fields{
			"messages-to-provide-evidence-for": slice.Map(msgsToProvideEvidenceFor, func(msg chain.MessageWithSignatures) uint64 {
				return msg.ID
			}),
		})
		logger.Info("providing evidence for messages")
		err := p.ProvideEvidence(ctx, queueName, msgsToProvideEvidenceFor)
		if err != nil {
			logger.WithError(err).Error("error providing evidence for messages")
			return err
		}
	}
	return nil
}

func (r *Relayer) Process(ctx context.Context, processors []chain.Processor) error {
	var processErrors whoops.Group

	if len(processors) == 0 {
		return nil
	}

	// todo randomise
	for _, p := range processors {

		// todo randomise
		for _, queueName := range p.SupportedQueues() {
			logger := log.WithFields(log.Fields{
				"queue-name": queueName,
			})

			messagesInQueue, err := r.palomaClient.QueryMessagesInQueue(ctx, queueName)
			logger.Debug("got ", len(messagesInQueue), " messages from ", queueName)
			if err != nil {
				logger.WithError(err).Error("couldn't get messages to relay")
				return err
			}

			messagesForSigning, err := r.palomaClient.QueryMessagesForSigning(ctx, queueName)
			if err != nil {
				logger.Error("failed getting messages to sign")
				return err
			}

			messagesForEvidence := messagesInQueue

			// only pick 5 messages to relay at a time.  This protects against timeouts
			messagesForRelaying := messagesInQueue
			if len(messagesForRelaying) > 5 {
				messagesForRelaying = messagesForRelaying[:5]
			}

			err = r.SignMessages(ctx, queueName, messagesForSigning, p)
			if err != nil {
				logger.Error("failed signing messages")
				processErrors.Add(err)
			}

			err = r.ProvideEvidenceForMessages(ctx, queueName, messagesForEvidence, p)
			if err != nil {
				logger.Error("failed providing evidence for messages")
				processErrors.Add(err)
			}

			err = r.RelayMessages(ctx, queueName, messagesForRelaying, p)
			if err != nil {
				logger.Error("failed relaying messages")
				processErrors.Add(err)
			}
		}
	}

	if processErrors.Err() {
		return processErrors
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
