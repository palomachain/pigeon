package relayer

import (
	"context"

	"github.com/palomachain/sparrow/chain"
)

func (r Relayer) Process(ctx context.Context) error {
	var processors []chain.Processor

	for _, p := range processors {
		for _, queueName := range p.SupportedQueues() {
			queuedMessages, err := r.palomaClient.QueryMessagesForSigning(ctx, r.validatorAddress, queueName)

			if err != nil {
				return err
			}
			p.SignMessages(ctx, queueName, queuedMessages...)
		}
	}
}
