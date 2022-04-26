package relayer

import (
	"context"

	"github.com/volumefi/conductor/client/paloma"
	cronchaintypes "github.com/volumefi/conductor/types/cronchain"
	"github.com/volumefi/utils/signing"
)

// THIS IS WIP AND WILL CHANGE!
// signMessagesForExecution takes messages from a given list of queueTypeNames that require a signature.
// It then signs each message and adds the signature into a list of signatures to be sent all at once
// over to the cronchain.
func (r *Relayer) signMessagesForExecution(ctx context.Context, queueTypeNames ...string) error {
	var broadcastMessageSignatures []paloma.BroadcastMessageSignatureIn
	for _, queueTypeName := range queueTypeNames {
		if err := ctx.Err(); err != nil {
			return err
		}
		// fetch messages that need to be signed
		msgs, err := paloma.QueryMessagesForSigning[*cronchaintypes.SignSmartContractExecute](
			ctx,
			r.palomaClient,
			// TODO: take the address from the keyring
			"VALIDATOR ADDRESS",
			queueTypeName,
		)
		if err != nil {
			return err
		}

		for _, msg := range msgs {
			// do the actual signing
			signBytes, _, err := signing.SignBytes(
				signing.KeyringSigner(r.palomaClient.Keyring(), r.palomaClient.L.Config.Key),
				signing.SerializeFnc(signing.JsonDeterministicEncoding),
				msg.Msg,
				msg.Nonce,
			)
			if err != nil {
				return err
			}
			broadcastMessageSignatures = append(
				broadcastMessageSignatures,
				paloma.BroadcastMessageSignatureIn{
					ID:            msg.ID,
					QueueTypeName: queueTypeName,
					Signature:     signBytes,
				},
			)
		}
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	return r.palomaClient.BroadcastMessageSignatures(ctx, broadcastMessageSignatures...)
}
