package relayer

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/palomachain/sparrow/chain/paloma"
	"github.com/palomachain/utils/signing"
	"github.com/vizualni/whoops"
)

// signMessagesForExecution takes messages from a given list of queue types and fetches and signes messages for every queue type provided.
// It is all aggregated and sent over in a single TX back to Paloma.
func (r *Relayer) signMessagesForExecution(ctx context.Context, queueNames ...string) error {
	var broadcastMessageSignatures []paloma.BroadcastMessageSignatureIn

	err := whoops.Try(func() {
		for _, queueName := range queueNames {
			whoops.Assert(ctx.Err())

			signingKeyAddress := whoops.Must(sdk.AccAddressFromBech32(r.signingKeyAddress))
			valAddress := whoops.Must(sdk.ValAddressFromBech32(r.validatorAddress))

			broadcastMessageSignatures = append(
				broadcastMessageSignatures,
				whoops.Must(
					signMessagesForExecution(
						ctx,
						r.palomaClient,
						signing.KeyringSignerByAddress(
							r.palomaClient.Keyring(),
							signingKeyAddress,
						),
						r.attestExecutor,
						valAddress,
						queueName,
					),
				)...,
			)
		}
	})

	if err != nil {
		return err
	}

	return r.palomaClient.BroadcastMessageSignatures(ctx, broadcastMessageSignatures...)
}

func signMessagesForExecution(
	ctx context.Context,
	client paloma.Client,
	signer signing.Signer,
	att attestExecutor,
	valAddress sdk.ValAddress,
	queueTypeName string,
) ([]paloma.BroadcastMessageSignatureIn, error) {
	var broadcastMessageSignatures []paloma.BroadcastMessageSignatureIn
	// fetch messages that need to be signed
	msgs, err := client.QueryMessagesForSigning(
		ctx,
		valAddress,
		queueTypeName,
	)
	if err != nil {
		return nil, err
	}

	for _, msg := range msgs {
		var extraData []byte

		// check if this is something that requires attestation
		evidence, err := att.Execute(ctx, queueTypeName, msg.Msg)
		if err != nil {
			return nil, err
		}

		if evidence != nil {
			extraData, err = evidence.Bytes()
			if err != nil {
				return nil, err
			}
		}

		// do the actual signing
		signBytes, _, err := signing.SignBytes(
			signer,
			signing.JsonDeterministicEncoding(),
			msg.Msg,
			msg.Nonce,
			extraData,
		)
		if err != nil {
			return nil, err
		}

		broadcastMessageSignatures = append(
			broadcastMessageSignatures,
			paloma.BroadcastMessageSignatureIn{
				ID:            msg.ID,
				QueueTypeName: queueTypeName,
				Signature:     signBytes,
				ExtraData:     extraData,
			},
		)
	}

	return broadcastMessageSignatures, nil
}
