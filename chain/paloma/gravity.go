package paloma

import (
	"context"
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/grpc"
	gravity "github.com/palomachain/paloma/x/gravity/types"
	"github.com/palomachain/pigeon/chain"
)

func (c *Client) GravityQueryLastUnsignedBatch(ctx context.Context, chainReferenceID string) ([]gravity.OutgoingTxBatch, error) {
	return gravityQueryLastUnsignedBatch(ctx, c.GRPCClient, c.creator, chainReferenceID)
}

func (c *Client) GravityConfirmBatches(ctx context.Context, signatures ...chain.SignedGravityOutgoingTxBatch) error {
	return gravityConfirmBatch(ctx, c.MessageSender, c.creator, signatures...)
}

func (c *Client) GravityQueryBatchesForRelaying(ctx context.Context, chainReferenceID string) ([]chain.GravityBatchWithSignatures, error) {
	return gravityQueryBatchesForRelaying(ctx, c.GRPCClient, c.valAddr, chainReferenceID)
}

func gravityConfirmBatch(
	ctx context.Context,
	ms MessageSender,
	creator string,
	signedBatches ...chain.SignedGravityOutgoingTxBatch,
) error {
	if len(signedBatches) == 0 {
		return nil
	}
	for _, signedBatch := range signedBatches {
		msg := &gravity.MsgConfirmBatch{
			Nonce:         signedBatch.BatchNonce,
			TokenContract: signedBatch.TokenContract,
			EthSigner:     signedBatch.SignedByAddress,
			Orchestrator:  creator,
			Signature:     hex.EncodeToString(signedBatch.Signature),
		}
		_, err := ms.SendMsg(ctx, msg, "")
		return err

	}
	return nil
}

func gravityQueryLastUnsignedBatch(ctx context.Context, grpcClient grpc.ClientConn, address string, chainReferenceID string) ([]gravity.OutgoingTxBatch, error) {
	qc := gravity.NewQueryClient(grpcClient)
	batches, err := qc.LastPendingBatchRequestByAddr(ctx, &gravity.QueryLastPendingBatchRequestByAddrRequest{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	filtered := make([]gravity.OutgoingTxBatch, 0, len(batches.Batch))
	for _, v := range batches.Batch {
		if v.GetChainReferenceID() == chainReferenceID {
			filtered = append(filtered, v)
		}
	}

	return filtered, nil
}

func gravityQueryBatchesForRelaying(ctx context.Context, grpcClient grpc.ClientConn, address sdk.ValAddress, chainReferenceID string) ([]chain.GravityBatchWithSignatures, error) {
	qc := gravity.NewQueryClient(grpcClient)

	// Get batches
	req := &gravity.QueryOutgoingTxBatchesRequest{
		ChainReferenceId: chainReferenceID,
		Assignee:         address.String(),
	}
	batches, err := qc.OutgoingTxBatches(ctx, req)
	if err != nil {
		return nil, err
	}

	batchesWithSignatures := make([]chain.GravityBatchWithSignatures, len(batches.Batches))
	for i, batch := range batches.Batches {
		confirms, err := qc.BatchConfirms(ctx, &gravity.QueryBatchConfirmsRequest{
			Nonce:           batch.BatchNonce,
			ContractAddress: batch.TokenContract,
		})
		if err != nil {
			return nil, err
		}

		var signatures []chain.ValidatorSignature
		for _, confirm := range confirms.Confirms {
			signature, err := hex.DecodeString(confirm.Signature)
			if err != nil {
				return nil, err
			}
			signatures = append(signatures, chain.ValidatorSignature{
				Signature:       signature,
				SignedByAddress: confirm.EthSigner,
			})
		}
		batchesWithSignatures[i] = chain.GravityBatchWithSignatures{
			OutgoingTxBatch: batch,
			Signatures:      signatures,
		}

	}

	return batchesWithSignatures, nil
}
