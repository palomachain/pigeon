package relayer

import (
	"context"
	"sync"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/internal/liblog"
	"github.com/palomachain/pigeon/internal/queue"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) AttestMessages(ctx context.Context, _ sync.Locker) error {
	logger := liblog.WithContext(ctx)
	logger.Info("attester loop")
	if ctx.Err() != nil {
		logger.Info("exiting attester loop as context has ended")
		return ctx.Err()
	}

	err := r.buildProcessors(ctx, nil)
	if err != nil {
		return err
	}

	err = r.attestMessages(ctx, r.processors)

	return handleProcessError(ctx, err)
}

func (r *Relayer) attestMessages(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	// todo randomise
	for _, p := range processors {
		// todo randomise
		for _, queueName := range p.SupportedQueues() {
			logger := liblog.WithContext(ctx).WithFields(log.Fields{
				"queue-name": queueName,
				"action":     "attest",
			})

			messagesInQueue, err := r.palomaClient.QueryMessagesForAttesting(ctx, queueName)
			if err != nil {
				logger.WithError(err).Error("couldn't get messages to attest")
				if isFatal(err) {
					return err
				}
				// Move on to the next queue on the same chain
				continue
			}

			logger = logger.WithFields(log.Fields{
				"message-ids": slice.Map(messagesInQueue, func(msg chain.MessageWithSignatures) uint64 {
					return msg.ID
				}),
			})

			logger.Debug("got ", len(messagesInQueue), " messages from ", queueName)

			if len(messagesInQueue) > 0 {
				logger := logger.WithFields(log.Fields{
					"messages-to-attest": slice.Map(messagesInQueue, func(msg chain.MessageWithSignatures) uint64 {
						return msg.ID
					}),
				})
				logger.Info("attesting ", len(messagesInQueue), " messages")
				err := p.ProvideEvidence(ctx, queue.FromString(queueName), messagesInQueue)
				if err != nil {
					logger.WithError(err).Error("error attesting messages")
					if err := r.palomaClient.
						NewStatus().
						WithChainReferenceID(p.GetChainReferenceID()).
						WithArg("error", err.Error()).
						WithLog("error attesting messages").
						Error(ctx); err != nil {
						logger.WithError(err).Error("failed to send Paloma status update")
					}

					if isFatal(err) {
						return err
					}
				}
			}
		}
	}

	return nil
}
