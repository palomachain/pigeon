package relayer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/internal/queue"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

const cMsgCacheSyncInterval = time.Minute

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
	if err != nil {
		return handleProcessError(ctx, err)
	}

	err = r.syncMsgCacheWithPaloma(ctx)
	return handleProcessError(ctx, err)
}

func (r *Relayer) syncMsgCacheWithPaloma(ctx context.Context) error {
	if time.Now().UTC().Sub(r.msgCache.lastSync) < cMsgCacheSyncInterval {
		return nil
	}

	msgIDs := maps.Keys(r.msgCache.records)
	if err := r.palomaClient.NewStatus().
		WithArg("msg-ids", fmt.Sprintf("%v", msgIDs)).
		WithLog("Query inbox sync.").
		Info(ctx); err != nil {
		return err
	}

	r.msgCache.lastSync = time.Now().UTC()
	clear(r.msgCache.records)

	return nil
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

			for _, v := range messagesInQueue {
				r.msgCache.records[v.ID] = struct{}{}
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
