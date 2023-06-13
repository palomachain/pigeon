package relayer

import (
	"context"
	"sync"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) AttestMessages(ctx context.Context, locker sync.Locker) error {
	log.Info("attester loop")
	if ctx.Err() != nil {
		log.Info("exiting attester loop as context has ended")
		return ctx.Err()
	}

	err := r.buildProcessors(ctx, locker)
	if err != nil {
		return err
	}

	locker.Lock()
	err = r.attestMessages(ctx, r.processors)
	locker.Unlock()

	return handleProcessError(err)
}

func (r *Relayer) attestMessages(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	// todo randomise
	for _, p := range processors {
		// todo randomise
		for _, queueName := range p.SupportedQueues() {
			logger := log.WithFields(log.Fields{
				"queue-name": queueName,
				"action":     "attest",
			})

			messagesInQueue, err := r.palomaClient.QueryMessagesForAttesting(ctx, queueName)

			logger = logger.WithFields(log.Fields{
				"message-ids": slice.Map(messagesInQueue, func(msg chain.MessageWithSignatures) uint64 {
					return msg.ID
				}),
			})

			logger.Debug("got ", len(messagesInQueue), " messages from ", queueName)
			if err != nil {
				logger.WithError(err).Error("couldn't get messages to attest")
				return err
			}

			msgsToAttest := slice.Filter(messagesInQueue, func(msg chain.MessageWithSignatures) bool {
				return len(msg.PublicAccessData) > 0 || len(msg.ErrorData) > 0
			})

			if len(msgsToAttest) > 0 {
				logger := logger.WithFields(log.Fields{
					"messages-to-attest": slice.Map(msgsToAttest, func(msg chain.MessageWithSignatures) uint64 {
						return msg.ID
					}),
				})
				logger.Info("attesting ", len(msgsToAttest), " messages")
				err := p.ProvideEvidence(ctx, queueName, msgsToAttest)
				if err != nil {
					logger.WithError(err).Error("error attesting messages")
					return err
				}
			}
		}
	}

	return nil
}
