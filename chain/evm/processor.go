package evm

import (
	"context"

	"github.com/palomachain/sparrow/chain"
)

type Processor struct {
	c *Client
}

var _ chain.Processor = Processor{}

func (p Processor) SupportedQueues() []string {
	return []string{"a", "b", "c"}
}

func (p Processor) SignMessages(ctx context.Context, queueTypeName string, messages ...chain.QueuedMessage) {
	for _, msg := range messages {
		p.c.sign(ctx, msg.BytesToSign)
	}

}

func (p Processor) RelayMessages(ctx context.Context, queueTypeName string, messages ...chain.ConsensusReachedMsg) {
}
