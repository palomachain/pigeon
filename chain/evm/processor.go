package evm

import (
	"context"
	"fmt"
	"strings"

	"github.com/palomachain/sparrow/chain"
	"github.com/palomachain/sparrow/types/paloma/x/evm/types"
	"github.com/palomachain/sparrow/util/slice"
)

const (
	queueArbitraryLogic = "evm-arbitrary-smart-contract-call"
)

type Processor struct {
	c         Client
	chainType string
	chainID   string
}

func NewProcessor(c Client, chainID string) Processor {
	return Processor{
		c:         c,
		chainType: "EVM",
		chainID:   chainID,
	}
}

var _ chain.Processor = Processor{}

func (p Processor) SupportedQueues() []string {
	return slice.Map(
		[]string{
			queueArbitraryLogic,
		},
		func(q string) string {
			return fmt.Sprintf("%s:%s:%s", p.chainType, p.chainID, q)
		},
	)
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
	switch {
	case strings.HasSuffix(queueTypeName, queueArbitraryLogic):
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
	default:
		return chain.ErrProcessorDoesNotSupportThisQueue.Format(queueTypeName)
	}
}

func (p Processor) ExternalAccount() chain.ExternalAccount {
	return chain.ExternalAccount{
		ChainType: p.chainType,
		ChainID:   p.chainID,
		Address:   p.c.addr.Hex(),
		PubKey:    p.c.addr.Bytes(),
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
