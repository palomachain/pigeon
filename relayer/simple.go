package relayer

import (
	"context"
	"fmt"

	"github.com/palomachain/sparrow/chain"
	"github.com/palomachain/sparrow/chain/paloma"
)

func (r Relayer) Process(ctx context.Context) error {
	processors := r.processors

	for _, p := range processors {
		for _, queueName := range p.SupportedQueues() {
			// TODO: remove comments once signing is done on the paloma side.
			// queuedMessages, err := r.palomaClient.QueryMessagesForSigning(ctx, r.validatorAddress, queueName)

			// if err != nil {
			// 	return err
			// }
			// signedMessages, err := p.SignMessages(ctx, queueName, queuedMessages...)
			// if err != nil {
			// 	return err
			// }

			// if err = r.broadcastSignaturesAndProcessAttestation(ctx, queueName, signedMessages); err != nil {
			// 	fmt.Println(err)
			// 	return err
			// }

			relayCandidateMsgs, err := r.palomaClient.QueryMessagesInQueue(ctx, queueName)
			if err = p.ProcessMessages(ctx, queueName, relayCandidateMsgs); err != nil {
				fmt.Println("error processing a message", err)
				return err
			}
		}
	}

	return nil
}

func (r Relayer) broadcastSignaturesAndProcessAttestation(ctx context.Context, queueTypeName string, sigs []chain.SignedQueuedMessage) error {
	var broadcastMessageSignatures []paloma.BroadcastMessageSignatureIn
	for _, sig := range sigs {
		var extraData []byte

		// check if this is something that requires attestation
		evidence, err := r.attestExecutor.Execute(ctx, queueTypeName, sig.Msg)
		if err != nil {
			return err
		}

		if evidence != nil {
			extraData, err = evidence.Bytes()
			if err != nil {
				return err
			}
		}

		broadcastMessageSignatures = append(broadcastMessageSignatures, paloma.BroadcastMessageSignatureIn{
			ID:            sig.ID,
			QueueTypeName: queueTypeName,
			Signature:     sig.Signature,
			ExtraData:     extraData,
		})
	}

	return r.palomaClient.BroadcastMessageSignatures(ctx)
}
