package paloma

import (
	"context"
	"strings"

	"github.com/VolumeFi/whoops"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/grpc"
	consensus "github.com/palomachain/paloma/x/consensus/types"
	evm "github.com/palomachain/paloma/x/evm/types"
	skyway "github.com/palomachain/paloma/x/skyway/types"
	valset "github.com/palomachain/paloma/x/valset/types"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/internal/liblog"
	"github.com/sirupsen/logrus"
)

// QueryMessagesForSigning returns a list of messages from a given queueTypeName that
// need to be signed by the provided validator given the valAddress.
func (c *Client) QueryMessagesForSigning(
	ctx context.Context,
	queueTypeName string,
) ([]chain.QueuedMessage, error) {
	return queryMessagesForSigning(
		ctx,
		c.GRPCClient,
		c.Unpacker,
		c.valAddr,
		queueTypeName,
	)
}

// QueryMessagesForAttesting returns all messages that are currently in the queue except those already attested for.
func (c *Client) QueryMessagesForAttesting(ctx context.Context, queueTypeName string) ([]chain.MessageWithSignatures, error) {
	return queryMessagesForAttesting(
		ctx,
		queueTypeName,
		c.valAddr,
		c.GRPCClient,
		c.Unpacker,
	)
}

// QueryMessagesForRelaying returns all messages that are currently in the queue.
func (c *Client) QueryMessagesForRelaying(ctx context.Context, queueTypeName string) ([]chain.MessageWithSignatures, error) {
	return queryMessagesForRelaying(
		ctx,
		queueTypeName,
		c.valAddr,
		c.GRPCClient,
		c.Unpacker,
	)
}

// QueryValidatorInfo returns info about the validator.
func (c *Client) QueryValidatorInfo(ctx context.Context) ([]*valset.ExternalChainInfo, error) {
	qc := valset.NewQueryClient(c.GRPCClient)
	valInfoRes, err := qc.ValidatorInfo(ctx, &valset.QueryValidatorInfoRequest{
		ValAddr: c.creatorValoper,
	})
	if err != nil {
		if strings.Contains(err.Error(), "item not found in store") {
			return nil, nil
		}
		return nil, err
	}

	return valInfoRes.ChainInfos, nil
}

// QueryGetSnapshotByID returns the snapshot by id. If the EventNonce is zero, then it returns the last snapshot.
func (c *Client) QueryGetSnapshotByID(ctx context.Context, id uint64) (*valset.Snapshot, error) {
	qc := valset.NewQueryClient(c.GRPCClient)
	snapshotRes, err := qc.GetSnapshotByID(ctx, &valset.QueryGetSnapshotByIDRequest{
		SnapshotId: id,
	})
	if err != nil {
		if strings.Contains(err.Error(), "item not found in store") {
			return nil, whoops.Enrich(
				chain.ErrNotFound,
				chain.EnrichedItemType.Val("snapshot"),
				chain.EnrichedID.Val(id),
			)
		}
		return nil, err
	}

	return snapshotRes.Snapshot, nil
}

func (c *Client) QueryGetLatestPublishedSnapshot(ctx context.Context, chainReferenceID string) (*valset.Snapshot, error) {
	qc := valset.NewQueryClient(c.GRPCClient)
	res, err := qc.GetLatestPublishedSnapshot(ctx, &valset.QueryGetLatestPublishedSnapshotRequest{
		ChainReferenceID: chainReferenceID,
	})
	if err != nil {
		return nil, err
	}

	return res.Snapshot, nil
}

func (c *Client) QueryLastObservedSkywayNonceByAddr(ctx context.Context, chainReferenceID string, orchestrator string) (uint64, error) {
	qc := skyway.NewQueryClient(c.GRPCClient)
	res, err := qc.LastObservedSkywayNonceByAddr(ctx, &skyway.QueryLastObservedSkywayNonceByAddrRequest{
		Address:          orchestrator,
		ChainReferenceId: chainReferenceID,
	})
	if err != nil {
		return 0, err
	}
	return res.Nonce, nil
}

func (c *Client) QueryBatchRequestByNonce(ctx context.Context, nonce uint64, contract string) (skyway.OutgoingTxBatch, error) {
	qc := skyway.NewQueryClient(c.GRPCClient)
	res, err := qc.BatchRequestByNonce(ctx, &skyway.QueryBatchRequestByNonceRequest{
		Nonce:           nonce,
		ContractAddress: contract,
	})
	if err != nil {
		return skyway.OutgoingTxBatch{}, err
	}
	return res.Batch, nil
}

func (c *Client) QueryGetEVMValsetByID(ctx context.Context, id uint64, chainReferenceID string) (*evm.Valset, error) {
	logger := liblog.WithContext(ctx)
	qc := evm.NewQueryClient(c.GRPCClient)
	valsetRes, err := qc.GetValsetByID(ctx, &evm.QueryGetValsetByIDRequest{
		ValsetID:         id,
		ChainReferenceID: chainReferenceID,
	})
	logger.WithFields(logrus.Fields{
		"valset-length":      len(valsetRes.Valset.Validators),
		"power-length":       len(valsetRes.Valset.Powers),
		"valset-id-out":      valsetRes.Valset.ValsetID,
		"valset-id-in":       id,
		"chain-reference-id": chainReferenceID,
	}).Debug("got valset by id")
	if err != nil {
		if strings.Contains(err.Error(), "item not found in store") {
			return nil, whoops.Enrich(
				chain.ErrNotFound,
				chain.EnrichedChainReferenceID.Val(chainReferenceID),
				chain.EnrichedID.Val(id),
				chain.EnrichedItemType.Val("valset"),
			)
		}
		return nil, err
	}

	return valsetRes.Valset, nil
}

func (c *Client) QueryGetEVMChainInfos(ctx context.Context) ([]*evm.ChainInfo, error) {
	qc := evm.NewQueryClient(c.GRPCClient)
	chainInfosRes, err := qc.ChainsInfos(ctx, &evm.QueryChainsInfosRequest{})
	if err != nil {
		return nil, err
	}

	return chainInfosRes.ChainsInfos, nil
}

func (c *Client) QueryGetValidatorAliveUntilBlockHeight(ctx context.Context) (int64, error) {
	qc := valset.NewQueryClient(c.GRPCClient)
	aliveUntilRes, err := qc.GetValidatorAliveUntil(ctx, &valset.QueryGetValidatorAliveUntilRequest{
		ValAddress: c.valAddr,
	})
	if err != nil {
		return 0, err
	}

	return aliveUntilRes.AliveUntilBlockHeight, nil
}

func queryMessagesForSigning(
	ctx context.Context,
	c grpc.ClientConn,
	anyunpacker codectypes.AnyUnpacker,
	valAddress sdk.ValAddress,
	queueTypeName string,
) ([]chain.QueuedMessage, error) {
	qc := consensus.NewQueryClient(c)
	msgs, err := qc.QueuedMessagesForSigning(ctx, &consensus.QueryQueuedMessagesForSigningRequest{
		ValAddress:    valAddress,
		QueueTypeName: queueTypeName,
	})
	if err != nil {
		return nil, err
	}
	res := make([]chain.QueuedMessage, len(msgs.GetMessageToSign()))
	for i, msg := range msgs.GetMessageToSign() {
		var ptr consensus.ConsensusMsg
		err := anyunpacker.UnpackAny(msg.GetMsg(), &ptr)
		if err != nil {
			return nil, err
		}
		res[i] = chain.QueuedMessage{
			ID:          msg.GetId(),
			Nonce:       msg.GetNonce(),
			BytesToSign: msg.GetBytesToSign(),
			Msg:         ptr,
		}
	}

	return res, nil
}

func queryMessagesForRelaying(
	ctx context.Context,
	queueTypeName string,
	valAddress sdk.ValAddress,
	c grpc.ClientConn,
	anyunpacker codectypes.AnyUnpacker,
) ([]chain.MessageWithSignatures, error) {
	qc := consensus.NewQueryClient(c)
	msgs, err := qc.QueuedMessagesForRelaying(ctx, &consensus.QueryQueuedMessagesForRelayingRequest{
		QueueTypeName: queueTypeName,
		ValAddress:    valAddress,
	})
	if err != nil {
		return nil, err
	}

	msgsWithSig := make([]chain.MessageWithSignatures, len(msgs.Messages))

	for i, msg := range msgs.Messages {
		valSigs := make([]chain.ValidatorSignature, len(msg.SignData))
		for j, vs := range msg.SignData {
			valSigs[j] = chain.ValidatorSignature{
				Signature:       vs.GetSignature(),
				SignedByAddress: vs.GetExternalAccountAddress(),
			}
		}
		var ptr consensus.ConsensusMsg
		err := anyunpacker.UnpackAny(msg.GetMsg(), &ptr)
		if err != nil {
			return nil, err
		}
		msgsWithSig[i] = chain.MessageWithSignatures{
			QueuedMessage: chain.QueuedMessage{
				ID:               msg.Id,
				Nonce:            msg.Nonce,
				Msg:              ptr,
				BytesToSign:      msg.GetBytesToSign(),
				PublicAccessData: msg.GetPublicAccessData(),
				ErrorData:        msg.GetErrorData(),
			},
			Signatures: valSigs,
		}
	}
	return msgsWithSig, err
}

func queryMessagesForAttesting(
	ctx context.Context,
	queueTypeName string,
	valAddress sdk.ValAddress,
	c grpc.ClientConn,
	anyunpacker codectypes.AnyUnpacker,
) ([]chain.MessageWithSignatures, error) {
	qc := consensus.NewQueryClient(c)
	msgs, err := qc.QueuedMessagesForAttesting(ctx, &consensus.QueryQueuedMessagesForAttestingRequest{
		QueueTypeName: queueTypeName,
		ValAddress:    valAddress,
	})
	if err != nil {
		return nil, err
	}

	msgsWithSig := make([]chain.MessageWithSignatures, len(msgs.Messages))
	for i, msg := range msgs.Messages {
		valSigs := make([]chain.ValidatorSignature, len(msg.SignData))
		for j, vs := range msg.SignData {
			valSigs[j] = chain.ValidatorSignature{
				// ValAddress:      vs.GetValAddress(),
				Signature:       vs.GetSignature(),
				SignedByAddress: vs.GetExternalAccountAddress(),
				// PublicKey:       vs.GetPublicKey(),
			}
		}
		var ptr consensus.ConsensusMsg
		err := anyunpacker.UnpackAny(msg.GetMsg(), &ptr)
		if err != nil {
			return nil, err
		}
		msgsWithSig[i] = chain.MessageWithSignatures{
			QueuedMessage: chain.QueuedMessage{
				ID:               msg.Id,
				Nonce:            msg.Nonce,
				Msg:              ptr,
				BytesToSign:      msg.GetBytesToSign(),
				PublicAccessData: msg.GetPublicAccessData(),
				ErrorData:        msg.GetErrorData(),
			},
			Signatures: valSigs,
		}
	}
	return msgsWithSig, err
}
