package relayer

import (
	"context"
	"sync"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) RelayMessages(ctx context.Context, locker sync.Locker) error {
	log.Info("relayer loop")
	if ctx.Err() != nil {
		log.Info("exiting relayer loop as context has ended")
		return ctx.Err()
	}

	err := r.buildProcessors(ctx, locker)
	if err != nil {
		return err
	}

	locker.Lock()
	err = r.relayMessages(ctx, r.processors)
	locker.Unlock()

	return handleProcessError(err)
}

func (r *Relayer) getMsgOffsetData(ctx context.Context) (int, int, error) {
	snapshot, err := r.palomaClient.QueryGetSnapshotByID(ctx, 0)
	if err != nil {
		return 0, 0, err
	}

	numValidators := len(snapshot.Validators)
	for i, validator := range snapshot.Validators {
		if r.palomaClient.GetValidatorAddress().String() == validator.Address.String() {
			return numValidators, i, nil
		}
	}

	// If we made it here, somehow our validator address wasn't in the snapshot
	return 0, 0, chain.ErrNotValidator
}

func (r *Relayer) relayMessages(ctx context.Context, processors []chain.Processor) error {
	if len(processors) == 0 {
		return nil
	}

	// Get the offset for this validator so that we know which messages to relay
	numValidators, myMsgOffset, err := r.getMsgOffsetData(ctx)
	if err != nil {
		return err
	}

	// todo randomise
	for _, p := range processors {
		// todo randomise
		for _, queueName := range p.SupportedQueues() {
			logger := log.WithFields(log.Fields{
				"queue-name": queueName,
				"action":     "relay",
			})

			messagesInQueue, err := r.palomaClient.QueryMessagesInQueue(ctx, queueName)

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

			relayCandidateMsgs := slice.Filter(
				messagesInQueue,
				func(msg chain.MessageWithSignatures) bool {
					return len(msg.PublicAccessData) == 0
				},
				func(msg chain.MessageWithSignatures) bool {
					return (int(msg.ID) % numValidators) == myMsgOffset
				},
			)

			if len(relayCandidateMsgs) > 0 {
				logger := logger.WithFields(log.Fields{
					"messages-to-relay": slice.Map(relayCandidateMsgs, func(msg chain.MessageWithSignatures) uint64 {
						return msg.ID
					}),
				})
				logger.Info("relaying ", len(relayCandidateMsgs), " messages")
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
		}
	}
	return nil
}
