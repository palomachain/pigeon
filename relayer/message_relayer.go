package relayer

import (
	"context"
	"sync"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/internal/queue"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) RelayMessages(ctx context.Context, _ sync.Locker) error {
	log.Info("relayer loop")
	if ctx.Err() != nil {
		log.Info("exiting relayer loop as context has ended")
		return ctx.Err()
	}

	err := r.buildProcessors(ctx, nil)
	if err != nil {
		return err
	}

	err = r.relayMessages(ctx, r.processors)

	return handleProcessError(ctx, err)
}

func (r *Relayer) relayMessages(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	// todo randomise
	for _, p := range processors {
		// todo randomise
		for _, queueName := range p.SupportedQueues() {
			logger := log.WithFields(log.Fields{
				"queue-name": queueName,
				"action":     "relay",
			})

			messagesInQueue, err := r.palomaClient.QueryMessagesForRelaying(ctx, queueName)

			logger = logger.WithFields(log.Fields{
				"message-ids": slice.Map(messagesInQueue, func(msg chain.MessageWithSignatures) uint64 {
					return msg.ID
				}),
			})

			logger.Debug("got ", len(messagesInQueue), " messages from ", queueName)
			if err != nil {
				logger.WithError(err).Error("couldn't get messages to relay")
				return err
			}

			if len(messagesInQueue) > 0 {
				logger := logger.WithFields(log.Fields{
					"messages-to-relay": slice.Map(messagesInQueue, func(msg chain.MessageWithSignatures) uint64 {
						return msg.ID
					}),
				})
				logger.Info("relaying ", len(messagesInQueue), " messages")
				err := p.ProcessMessages(ctx, queue.FromString(queueName), messagesInQueue)
				if err != nil {
					logger.WithFields(log.Fields{
						"err":        err,
						"queue-name": queueName,
						"messages-to-relay": slice.Map(messagesInQueue, func(msg chain.MessageWithSignatures) uint64 {
							return msg.ID
						}),
					}).Error("error relaying messages")
					return err
				}
			}
		}
	}
	return nil
}
