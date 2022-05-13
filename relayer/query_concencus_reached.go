package relayer

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/palomachain/sparrow/client/paloma"
	consensus "github.com/palomachain/sparrow/types/paloma/x/consensus/types"
	"github.com/palomachain/utils/signing"
	"github.com/vizualni/whoops"
)

// signMessagesForExecution takes messages from a given list of queue types and
// fetches and signes messages for every queue type provided.  It is all
// aggregated and sent over in a single TX back to Paloma.
func (r *Relayer) relayConsensusReachedMessages(ctx context.Context, qtss ...queueTypeSignerer) error {
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

func (c consensusMessageQueueType[T]) relayConsensusReachedMessages(
	ctx context.Context,
	r *Relayer,
) ([]paloma.BroadcastMessageSignatureIn, error) {
	return signMessagesForExecution[T](
		ctx,
		r.palomaClient,
		signing.KeyringSignerByAddress(
			r.palomaClient.Keyring(),
			whoops.Must(sdk.AccAddressFromBech32(r.signingKeyAddress)),
		),
		r.signingKeyAddress,
		c.queue(),
	)
}

func relayConsensusReachedMessages[T consensus.Signable](
	ctx context.Context,
	client paloma.Client,
	signer signing.Signer,
	signingKeyAddress string,
	queueTypeName string,
) ([]paloma.ConsensusReachedMsg[T], error) {
	var broadcastMessageSignatures []paloma.BroadcastMessageSignatureIn
	// fetch messages that need to be signed
	msgs, err := paloma.QueryConsensusReachedMessages[T](
		ctx,
		client,
		queueTypeName,
	)
	if err != nil {
		return nil, err
	}

	for _, msg := range msgs {
		// do the actual signing
		signBytes, _, err := signing.SignBytes(
			signer,
			signing.SerializeFnc(signing.JsonDeterministicEncoding),
			msg.Msg,
			msg.Nonce,
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
			},
		)
	}

	return broadcastMessageSignatures, nil
}
