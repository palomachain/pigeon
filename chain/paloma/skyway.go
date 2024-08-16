package paloma

import (
	"context"
	"encoding/hex"

	skyway "github.com/palomachain/paloma/x/skyway/types"
	"github.com/palomachain/pigeon/chain"
)

func (c *Client) SkywayQueryLastUnsignedBatch(ctx context.Context, chainReferenceID string) ([]skyway.OutgoingTxBatch, error) {
	// return skywayQueryLastUnsignedBatch(ctx, c.GRPCClient, c.creator, chainReferenceID)
	qc := skyway.NewQueryClient(c.GRPCClient)
	batches, err := qc.LastPendingBatchRequestByAddr(ctx, &skyway.QueryLastPendingBatchRequestByAddrRequest{
		Address: c.creator,
	})
	if err != nil {
		return nil, err
	}

	filtered := make([]skyway.OutgoingTxBatch, 0, len(batches.Batch))
	for _, v := range batches.Batch {
		if v.GetChainReferenceID() == chainReferenceID {
			filtered = append(filtered, v)
		}
	}

	return filtered, nil
}

func (c *Client) SkywayConfirmBatches(ctx context.Context, signatures ...chain.SignedSkywayOutgoingTxBatch) error {
	if len(signatures) == 0 {
		return nil
	}
	for _, signedBatch := range signatures {
		msg := &skyway.MsgConfirmBatch{
			Nonce:         signedBatch.BatchNonce,
			TokenContract: signedBatch.TokenContract,
			EthSigner:     signedBatch.SignedByAddress,
			Orchestrator:  c.creator,
			Signature:     hex.EncodeToString(signedBatch.Signature),
		}
		_, err := c.messageSender.SendMsg(ctx, msg, "", c.sendingOpts...)
		return err

	}
	return nil
}

func (c *Client) SkywayQueryBatchesForRelaying(ctx context.Context, chainReferenceID string) ([]chain.SkywayBatchWithSignatures, error) {
	// return skywayQueryBatchesForRelaying(ctx, c.GRPCClient, c.valAddr, chainReferenceID)
	qc := skyway.NewQueryClient(c.GRPCClient)

	// Get batches
	req := &skyway.QueryOutgoingTxBatchesRequest{
		ChainReferenceId: chainReferenceID,
		Assignee:         c.valAddr.String(),
	}
	batches, err := qc.OutgoingTxBatches(ctx, req)
	if err != nil {
		return nil, err
	}

	return skywayLoadConfirms(ctx, qc, batches.Batches)
}

func (c *Client) SkywayQueryLastPendingBatchForGasEstimation(ctx context.Context, chainReferenceID string) ([]chain.SkywayBatchWithSignatures, error) {
	qc := skyway.NewQueryClient(c.GRPCClient)
	res, err := qc.LastPendingBatchForGasEstimation(ctx, &skyway.QueryLastPendingBatchForGasEstimationRequest{
		Address:          c.valAddr,
		ChainReferenceId: chainReferenceID,
	})
	if err != nil {
		return nil, err
	}

	return skywayLoadConfirms(ctx, qc, res.Batch)
}

func (c *Client) SkywayEstimateBatchGas(ctx context.Context, estimates ...chain.EstimatedSkywayBatch) error {
	if len(estimates) == 0 {
		return nil
	}
	for _, v := range estimates {
		msg := &skyway.MsgEstimateBatchGas{
			Nonce:         v.BatchNonce,
			TokenContract: v.TokenContract,
			EthSigner:     v.EstimatedByAddress,
			Estimate:      v.Value,
		}
		_, err := c.messageSender.SendMsg(ctx, msg, "", c.sendingOpts...)
		if err != nil {
			return err
		}

	}
	return nil
}

// TODO: Combine with below method
func (c *Client) SendBatchSendToEVMClaim(ctx context.Context, claim skyway.MsgBatchSendToRemoteClaim) error {
	_, err := c.messageSender.SendMsg(ctx, &claim, "", c.sendingOpts...)
	return err
}

func (c *Client) SendSendToPalomaClaim(ctx context.Context, claim skyway.MsgSendToPalomaClaim) error {
	_, err := c.messageSender.SendMsg(ctx, &claim, "", c.sendingOpts...)
	return err
}

func (c *Client) SendLightNodeSaleClaim(ctx context.Context, claim skyway.MsgLightNodeSaleClaim) error {
	_, err := c.messageSender.SendMsg(ctx, &claim, "", c.sendingOpts...)
	return err
}

func skywayLoadConfirms(ctx context.Context, qc skyway.QueryClient, in []skyway.OutgoingTxBatch) ([]chain.SkywayBatchWithSignatures, error) {
	batchesWithSignatures := make([]chain.SkywayBatchWithSignatures, len(in))
	for i, batch := range in {
		confirms, err := qc.BatchConfirms(ctx, &skyway.QueryBatchConfirmsRequest{
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
		batchesWithSignatures[i] = chain.SkywayBatchWithSignatures{
			OutgoingTxBatch: batch,
			Signatures:      signatures,
		}

	}
	return batchesWithSignatures, nil
}
