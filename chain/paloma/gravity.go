package paloma

import (
	"context"
	"encoding/hex"

	gravity "github.com/palomachain/paloma/x/gravity/types"
	"github.com/palomachain/pigeon/chain"
)

func (c *Client) GravityQueryLastUnsignedBatch(ctx context.Context, chainReferenceID string) ([]gravity.OutgoingTxBatch, error) {
	// return gravityQueryLastUnsignedBatch(ctx, c.GRPCClient, c.creator, chainReferenceID)
	qc := gravity.NewQueryClient(c.GRPCClient)
	batches, err := qc.LastPendingBatchRequestByAddr(ctx, &gravity.QueryLastPendingBatchRequestByAddrRequest{
		Address: c.creator,
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

func (c *Client) GravityConfirmBatches(ctx context.Context, signatures ...chain.SignedGravityOutgoingTxBatch) error {
	if len(signatures) == 0 {
		return nil
	}
	for _, signedBatch := range signatures {
		msg := &gravity.MsgConfirmBatch{
			Nonce:         signedBatch.BatchNonce,
			TokenContract: signedBatch.TokenContract,
			EthSigner:     signedBatch.SignedByAddress,
			Orchestrator:  c.creator,
			Signature:     hex.EncodeToString(signedBatch.Signature),
		}
		_, err := c.MessageSender.SendMsg(ctx, msg, "", c.sendingOpts...)
		return err

	}
	return nil
}

func (c *Client) GravityQueryBatchesForRelaying(ctx context.Context, chainReferenceID string) ([]chain.GravityBatchWithSignatures, error) {
	// return gravityQueryBatchesForRelaying(ctx, c.GRPCClient, c.valAddr, chainReferenceID)
	qc := gravity.NewQueryClient(c.GRPCClient)

	// Get batches
	req := &gravity.QueryOutgoingTxBatchesRequest{
		ChainReferenceId: chainReferenceID,
		Assignee:         c.valAddr.String(),
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

// TODO: Combine with below method
func (c *Client) SendBatchSendToEVMClaim(ctx context.Context, claim gravity.MsgBatchSendToEthClaim) error {
	_, err := c.MessageSender.SendMsg(ctx, &claim, "", c.sendingOpts...)
	return err
}

func (c *Client) SendSendToPalomaClaim(ctx context.Context, claim gravity.MsgSendToPalomaClaim) error {
	_, err := c.MessageSender.SendMsg(ctx, &claim, "", c.sendingOpts...)
	return err
}
