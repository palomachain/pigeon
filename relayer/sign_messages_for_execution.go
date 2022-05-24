package relayer

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/palomachain/sparrow/chain/paloma"
	consensus "github.com/palomachain/sparrow/types/paloma/x/consensus/types"
	"github.com/palomachain/utils/signing"
	"github.com/vizualni/whoops"
)

type queueTypeSignerer interface {
	queryMessagesForSigning(ctx context.Context, r *Relayer) ([]paloma.BroadcastMessageSignatureIn, error)
}

// signMessagesForExecution takes messages from a given list of queue types and fetches and signes messages for every queue type provided.
// It is all aggregated and sent over in a single TX back to Paloma.
func (r *Relayer) signMessagesForExecution(ctx context.Context, qtss ...queueTypeSignerer) error {
	var broadcastMessageSignatures []paloma.BroadcastMessageSignatureIn

	err := whoops.Try(func() {
		for _, qts := range qtss {
			whoops.Assert(ctx.Err())

			broadcastMessageSignatures = append(
				broadcastMessageSignatures,
				whoops.Must(
					qts.queryMessagesForSigning(ctx, r),
				)...,
			)
		}
	})

	if err != nil {
		return err
	}

	return r.palomaClient.BroadcastMessageSignatures(ctx, broadcastMessageSignatures...)
}

func (c consensusMessageQueueType[T]) queryMessagesForSigning(
	ctx context.Context,
	r *Relayer,
) ([]paloma.BroadcastMessageSignatureIn, error) {
	signingKeyAddress := whoops.Must(sdk.AccAddressFromBech32(r.signingKeyAddress))
	valAddress := whoops.Must(sdk.ValAddressFromBech32(r.validatorAddress))
	return signMessagesForExecution[T](
		ctx,
		r.palomaClient,
		signing.KeyringSignerByAddress(
			r.palomaClient.Keyring(),
			signingKeyAddress,
		),
		r.attestExecutor,
		valAddress,
		c.queue(),
	)
}

func signMessagesForExecution[T consensus.Signable](
	ctx context.Context,
	client paloma.Client,
	signer signing.Signer,
	att attestExecutor,
	valAddress sdk.ValAddress,
	queueTypeName string,
) ([]paloma.BroadcastMessageSignatureIn, error) {
	var broadcastMessageSignatures []paloma.BroadcastMessageSignatureIn
	// fetch messages that need to be signed
	msgs, err := paloma.QueryMessagesForSigning[T](
		ctx,
		client,
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
