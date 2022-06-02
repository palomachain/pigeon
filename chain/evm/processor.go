package evm

import (
	"context"
	"fmt"

	"github.com/palomachain/sparrow/chain"
	"github.com/palomachain/sparrow/types/paloma/x/evm/types"
)

const (
	queueArbitraryLogic = "evm-arbitrary-smart-contract-call"
)

type Processor struct {
	c Client
}

func NewProcessor(c Client) Processor {
	return Processor{c}
}

var _ chain.Processor = Processor{}

func (p Processor) SupportedQueues() []string {
	return []string{queueArbitraryLogic}
}

func (p Processor) SignMessages(ctx context.Context, queueTypeName string, messages ...chain.QueuedMessage) ([]chain.SignedQueuedMessage, error) {
	signed := []chain.SignedQueuedMessage{}
	for _, msg := range messages {
		sig, err := p.c.sign(ctx, msg.BytesToSign)
		if err != nil {
			return nil, err
		}
		signed = append(signed, chain.SignedQueuedMessage{
			QueuedMessage: msg,
			Signature:     sig,
		})
	}

	return signed, nil
}

func (p Processor) ProcessMessages(ctx context.Context, queueTypeName string, msgs []chain.MessageWithSignatures) error {
	// TODO: check for signatures
	switch queueTypeName {
	case queueArbitraryLogic:
		return p.processArbitraryLogic(
			ctx,
			queueTypeName,
			typeMapSlice(
				msgs,
				func(msg chain.MessageWithSignatures) *types.ArbitrarySmartContractCall {
					return msg.Msg.(*types.ArbitrarySmartContractCall)
				},
			),
			typeMapSlice(
				msgs,
				func(msg chain.MessageWithSignatures) uint64 {
					return msg.ID
				},
			),
		)
		return nil
	default:
		return chain.ErrProcessorDoesNotSupportThisQueue.Format(queueTypeName)
	}
}

// TODO don't use types.ArbitrarySmartContractCall
func (p Processor) processArbitraryLogic(ctx context.Context, queueTypeName string, msgs []*types.ArbitrarySmartContractCall, ids []uint64) error {
	for i, msg := range msgs {
		err := p.c.executeArbitraryMessage(ctx, msg)
		if err != nil {
			return err
		}
		fmt.Println("THIS IS TEMPORARY ONLY: DELETING JOB FROM QUEUE GIVEN THAT IT WAS SENT")
		err = p.c.paloma.DeleteJob(ctx, queueTypeName, ids[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func typeMapSlice[A any, B any](slice []A, fnc func(A) B) (res []B) {
	for _, item := range slice {
		res = append(res, fnc(item))
	}
	return
}
