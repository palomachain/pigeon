package relayer

import (
	"context"
	"sync"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/internal/queue"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) EstimateMessages(ctx context.Context, _ sync.Locker) error {
	log.Info("estimator loop")
	if ctx.Err() != nil {
		log.Info("exiting estimator loop as context has ended")
		return ctx.Err()
	}

	err := r.buildProcessors(ctx, nil)
	if err != nil {
		return err
	}

	err = r.estimateMessages(ctx, r.processors)
	if err != nil {
		return handleProcessError(ctx, err)
	}

	return handleProcessError(ctx, err)
}

func (r *Relayer) estimateMessages(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	for _, p := range processors {
		for _, queueName := range p.SupportedQueues() {
			logger := log.WithFields(log.Fields{
				"queue-name": queueName,
				"action":     "estimate",
			})

			messagesInQueue, err := r.palomaClient.QueryMessagesForEstimating(ctx, queueName)

			logger = logger.WithFields(log.Fields{
				"message-ids": slice.Map(messagesInQueue, func(msg chain.MessageWithSignatures) uint64 {
					return msg.ID
				}),
			})

			logger.Debug("got ", len(messagesInQueue), " messages from ", queueName)
			if err != nil {
				logger.WithError(err).Error("couldn't get messages to estimate")
				return err
			}

			if len(messagesInQueue) > 0 {
				logger := logger.WithFields(log.Fields{
					"messages-to-estimate": slice.Map(messagesInQueue, func(msg chain.MessageWithSignatures) uint64 {
						return msg.ID
					}),
				})
				logger.Info("estimating ", len(messagesInQueue), " messages")
				estimates, err := p.EstimateMessages(ctx, queue.FromString(queueName), messagesInQueue)
				if err != nil {
					logger.WithError(err).Error("error estimating messages")
					return err
				}

				filteredEstimates := make([]chain.MessageWithEstimate, 0, len(estimates))
				for _, v := range estimates {
					logger.
						WithField("message-id", v.ID).
						WithField("estimate", v.Estimate).
						Info("estimated message")
					if v.Estimate < 1 {
						logger.
							WithField("message-id", v.ID).
							WithField("estimate", v.Estimate).
							Warn("Received an estimate of 0 or less, skipping")
						continue
					}
					filteredEstimates = append(filteredEstimates, v)
				}

				if len(filteredEstimates) == 0 {
					logger.Info("No valid estimates, skipping")
					continue
				}

				err = r.palomaClient.AddMessagesGasEstimate(ctx, queueName, filteredEstimates...)
				if err != nil {
					logger.WithError(err).Error("failed to send estimates to Paloma")
					return err
				}
			}

		}
	}
	return nil
}
