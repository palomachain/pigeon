package conductor

import (
	"context"

	"github.com/volumefi/conductor/client/cronchain"
	cronchaintypes "github.com/volumefi/conductor/types/cronchain"
)

type cronchainerClienter interface {
}

type relayer struct {
	cronchain cronchain.Client
	terra     any
}

// signMessagesForExecution takes messages from a given list of queueTypeNames that require a signature.
// It then signs each message and adds the signature into a list of signatures to be sent all at once
// over to the cronchain.
func (r relayer) signMessagesForExecution(ctx context.Context, queueTypeNames ...string) error {
	var broadcastMessageSignatures []cronchain.BroadcastMessageSignatureIn
	for _, queueTypeName := range queueTypeNames {
		if err := ctx.Err(); err != nil {
			return err
		}
		// fetch messages that need to be signed
		msgs, err := cronchain.QueryMessagesForSigning[*cronchaintypes.SignSmartContractExecute](
			ctx,
			r.cronchain,
			// TODO: take the address from the keyring
			"VALIDATOR ADDRESS",
			queueTypeName,
		)
		if err != nil {
			return err
		}

		for _, msg := range msgs {
			// do the actual signing
			signBytes, _, err := signBytes(
				keyringSigner(r.cronchain.Keyring(), r.cronchain.L.Config.Key),
				serializeFnc(jsonDeterministicEncoding),
				msg.Msg,
				msg.Nonce,
			)
			if err != nil {
				return err
			}
			broadcastMessageSignatures = append(broadcastMessageSignatures, cronchain.BroadcastMessageSignatureIn{
				ID:            msg.ID,
				QueueTypeName: queueTypeName,
				Signature:     signBytes,
			})
		}
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	return r.cronchain.BroadcastMessageSignatures(ctx, broadcastMessageSignatures...)
}
