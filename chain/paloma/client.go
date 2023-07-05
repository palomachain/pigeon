package paloma

import (
	"context"
	"strings"
	"time"

	"github.com/VolumeFi/whoops"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	consensus "github.com/palomachain/paloma/x/consensus/types"
	evm "github.com/palomachain/paloma/x/evm/types"
	valset "github.com/palomachain/paloma/x/valset/types"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
	"github.com/strangelove-ventures/lens/client/query"
)

type ResultStatus = coretypes.ResultStatus

//go:generate mockery --name=MessageSender
type MessageSender interface {
	SendMsg(ctx context.Context, msg sdk.Msg, memo string) (*sdk.TxResponse, error)
}

type Client struct {
	L            *chain.LensClient
	PalomaConfig config.Paloma

	GRPCClient grpc.ClientConn

	MessageSender MessageSender

	creator        string
	creatorValoper string
	valAddr        sdk.ValAddress
}

func (c *Client) Init() {
	c.creator = getCreator(*c)
	c.creatorValoper = getCreatorAsValoper(*c)
	c.valAddr = sdk.ValAddress(getMainAddress(*c).Bytes())
}

// QueryMessagesForSigning returns a list of messages from a given queueTypeName that
// need to be signed by the provided validator given the valAddress.
func (c Client) QueryMessagesForSigning(
	ctx context.Context,
	queueTypeName string,
) ([]chain.QueuedMessage, error) {
	return queryMessagesForSigning(
		ctx,
		c.GRPCClient,
		c.L.Codec.Marshaler,
		c.valAddr,
		queueTypeName,
	)
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
	res := []chain.QueuedMessage{}
	for _, msg := range msgs.GetMessageToSign() {
		var ptr consensus.ConsensusMsg
		err := anyunpacker.UnpackAny(msg.GetMsg(), &ptr)
		if err != nil {
			return nil, err
		}
		res = append(res, chain.QueuedMessage{
			ID:          msg.GetId(),
			Nonce:       msg.GetNonce(),
			BytesToSign: msg.GetBytesToSign(),
			Msg:         ptr,
		})
	}

	return res, nil
}

// QueryMessagesForAttesting returns all messages that are currently in the queue except those already attested for.
func (c Client) QueryMessagesForAttesting(ctx context.Context, queueTypeName string) ([]chain.MessageWithSignatures, error) {
	return queryMessagesForAttesting(
		ctx,
		queueTypeName,
		c.valAddr,
		c.GRPCClient,
		c.L.Codec.Marshaler,
	)
}

// QueryMessagesForRelaying returns all messages that are currently in the queue.
func (c Client) QueryMessagesForRelaying(ctx context.Context, queueTypeName string) ([]chain.MessageWithSignatures, error) {
	return queryMessagesForRelaying(
		ctx,
		queueTypeName,
		c.valAddr,
		c.GRPCClient,
		c.L.Codec.Marshaler,
	)
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

	msgsWithSig := []chain.MessageWithSignatures{}
	for _, msg := range msgs.Messages {
		valSigs := []chain.ValidatorSignature{}
		for _, vs := range msg.SignData {
			valSigs = append(valSigs, chain.ValidatorSignature{
				ValAddress:      vs.GetValAddress(),
				Signature:       vs.GetSignature(),
				SignedByAddress: vs.GetExternalAccountAddress(),
				PublicKey:       vs.GetPublicKey(),
			})
		}
		var ptr consensus.ConsensusMsg
		err := anyunpacker.UnpackAny(msg.GetMsg(), &ptr)
		if err != nil {
			return nil, err
		}
		msgsWithSig = append(msgsWithSig, chain.MessageWithSignatures{
			QueuedMessage: chain.QueuedMessage{
				ID:               msg.Id,
				Nonce:            msg.Nonce,
				Msg:              ptr,
				BytesToSign:      msg.GetBytesToSign(),
				PublicAccessData: msg.GetPublicAccessData(),
				ErrorData:        msg.GetErrorData(),
			},
			Signatures: valSigs,
		})
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

	msgsWithSig := []chain.MessageWithSignatures{}
	for _, msg := range msgs.Messages {
		valSigs := []chain.ValidatorSignature{}
		for _, vs := range msg.SignData {
			valSigs = append(valSigs, chain.ValidatorSignature{
				ValAddress:      vs.GetValAddress(),
				Signature:       vs.GetSignature(),
				SignedByAddress: vs.GetExternalAccountAddress(),
				PublicKey:       vs.GetPublicKey(),
			})
		}
		var ptr consensus.ConsensusMsg
		err := anyunpacker.UnpackAny(msg.GetMsg(), &ptr)
		if err != nil {
			return nil, err
		}
		msgsWithSig = append(msgsWithSig, chain.MessageWithSignatures{
			QueuedMessage: chain.QueuedMessage{
				ID:               msg.Id,
				Nonce:            msg.Nonce,
				Msg:              ptr,
				BytesToSign:      msg.GetBytesToSign(),
				PublicAccessData: msg.GetPublicAccessData(),
				ErrorData:        msg.GetErrorData(),
			},
			Signatures: valSigs,
		})
	}
	return msgsWithSig, err
}

type BroadcastMessageSignatureIn struct {
	ID              uint64
	QueueTypeName   string
	Signature       []byte
	SignedByAddress string
}

// BroadcastMessageSignatures takes a list of signatures that need to be sent over to the chain.
// It build the message and sends it over.
func (c Client) BroadcastMessageSignatures(ctx context.Context, signatures ...BroadcastMessageSignatureIn) error {
	return broadcastMessageSignatures(ctx, c.MessageSender, c.creator, signatures...)
}

// QueryValidatorInfo returns info about the validator.
func (c Client) QueryValidatorInfo(ctx context.Context) ([]*valset.ExternalChainInfo, error) {
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

// QueryGetSnapshotByID returns the snapshot by id. If the ID is zero, then it returns the last snapshot.
func (c Client) QueryGetSnapshotByID(ctx context.Context, id uint64) (*valset.Snapshot, error) {
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

func (c Client) BlockHeight(ctx context.Context) (int64, error) {
	res, err := c.L.RPCClient.Status(ctx)
	if err != nil {
		return 0, err
	}

	if res.SyncInfo.CatchingUp {
		return 0, ErrNodeIsNotInSync
	}

	return res.SyncInfo.LatestBlockHeight, nil
}

func (c Client) QueryGetEVMValsetByID(ctx context.Context, id uint64, chainReferenceID string) (*evm.Valset, error) {
	qc := evm.NewQueryClient(c.GRPCClient)
	valsetRes, err := qc.GetValsetByID(ctx, &evm.QueryGetValsetByIDRequest{
		ValsetID:         id,
		ChainReferenceID: chainReferenceID,
	})
	log.WithFields(log.Fields{
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

// TODO: this should return all chain infos. Not the ones from EVM only.
func (c Client) QueryGetEVMChainInfos(ctx context.Context) ([]*evm.ChainInfo, error) {
	qc := evm.NewQueryClient(c.GRPCClient)
	chainInfosRes, err := qc.ChainsInfos(ctx, &evm.QueryChainsInfosRequest{})
	if err != nil {
		return nil, err
	}

	return chainInfosRes.ChainsInfos, nil
}

// TODO: this is only temporary for easier testing
func (c Client) DeleteJob(ctx context.Context, queueTypeName string, id uint64) error {
	_, err := c.MessageSender.SendMsg(ctx, &consensus.MsgDeleteJob{
		Creator:       c.creator,
		QueueTypeName: queueTypeName,
		MessageID:     id,
	}, "")
	return err
}

type ChainInfoIn struct {
	ChainType        string
	ChainReferenceID string
	AccAddress       string
	PubKey           []byte
}

// AddExternalChainInfo adds info about the external chain. It adds the chain's
// account addresses that the pigeon knows about.
func (c Client) AddExternalChainInfo(ctx context.Context, chainInfos ...ChainInfoIn) error {
	if len(chainInfos) == 0 {
		return nil
	}

	msg := &valset.MsgAddExternalChainInfoForValidator{
		Creator: c.creator,
	}

	msg.ChainInfos = slice.Map(
		chainInfos,
		func(in ChainInfoIn) *valset.ExternalChainInfo {
			return &valset.ExternalChainInfo{
				ChainType:        in.ChainType,
				ChainReferenceID: in.ChainReferenceID,
				Address:          in.AccAddress,
				Pubkey:           in.PubKey,
			}
		},
	)

	_, err := c.MessageSender.SendMsg(ctx, msg, "")
	return err
}

func (c Client) AddMessageEvidence(ctx context.Context, queueTypeName string, messageID uint64, proof proto.Message) error {
	anyProof, err := codectypes.NewAnyWithValue(proof)
	if err != nil {
		return err
	}
	msg := &consensus.MsgAddEvidence{
		Creator:       c.creator,
		Proof:         anyProof,
		MessageID:     messageID,
		QueueTypeName: queueTypeName,
	}

	_, err = c.MessageSender.SendMsg(ctx, msg, "")
	return err
}

func (c Client) SetPublicAccessData(ctx context.Context, queueTypeName string, messageID uint64, data []byte) error {
	msg := &consensus.MsgSetPublicAccessData{
		Creator:       c.creator,
		Data:          data,
		MessageID:     messageID,
		QueueTypeName: queueTypeName,
	}

	_, err := c.MessageSender.SendMsg(ctx, msg, "")
	return err
}

func (c Client) SetErrorData(ctx context.Context, queueTypeName string, messageID uint64, data []byte) error {
	msg := &consensus.MsgSetErrorData{
		Creator:       c.creator,
		Data:          data,
		MessageID:     messageID,
		QueueTypeName: queueTypeName,
	}

	_, err := c.MessageSender.SendMsg(ctx, msg, "")
	return err
}

func (c Client) QueryGetValidatorAliveUntil(ctx context.Context) (time.Time, error) {
	qc := valset.NewQueryClient(c.GRPCClient)
	aliveUntilRes, err := qc.GetValidatorAliveUntil(ctx, &valset.QueryGetValidatorAliveUntilRequest{
		ValAddress: c.valAddr,
	})
	if err != nil {
		return time.Time{}, err
	}

	return aliveUntilRes.AliveUntil.UTC(), nil
}

func (c Client) KeepValidatorAlive(ctx context.Context) error {
	msg := &valset.MsgKeepAlive{
		Creator: c.creator,
	}

	_, err := c.MessageSender.SendMsg(ctx, msg, "")
	return err
}

func (c Client) lensQuery() *query.Query {
	return &query.Query{Client: &c.L.ChainClient, Options: query.DefaultOptions()}
}

func (c Client) Status(ctx context.Context) (*ResultStatus, error) {
	return c.lensQuery().Status()
}

func (c Client) PalomaStatus(ctx context.Context) error {
	res, err := c.Status(ctx)

	if IsPalomaDown(err) {
		return whoops.Wrap(ErrPalomaIsDown, err)
	}

	if err != nil {
		return err
	}
	_ = res
	return nil
}

func (c Client) GetValidator(ctx context.Context) (*stakingtypes.Validator, error) {
	res, err := c.lensQuery().Staking_Validator(c.GetValidatorAddress().String())
	if err != nil {
		return nil, err
	}
	return &res.Validator, nil
}

func broadcastMessageSignatures(
	ctx context.Context,
	ms MessageSender,
	creator string,
	signatures ...BroadcastMessageSignatureIn,
) error {
	if len(signatures) == 0 {
		return nil
	}
	var signedMessages []*consensus.ConsensusMessageSignature
	for _, sig := range signatures {
		signedMessages = append(signedMessages, &consensus.ConsensusMessageSignature{
			Id:              sig.ID,
			QueueTypeName:   sig.QueueTypeName,
			Signature:       sig.Signature,
			SignedByAddress: sig.SignedByAddress,
		})
	}
	msg := &consensus.MsgAddMessagesSignatures{
		Creator:        creator,
		SignedMessages: signedMessages,
	}
	_, err := ms.SendMsg(ctx, msg, "")
	return err
}

func (c Client) Keyring() keyring.Keyring {
	return c.L.Keybase
}

func (c Client) GetValidatorAddress() sdk.ValAddress {
	return c.valAddr
}

func getMainAddress(c Client) sdk.Address {
	key, err := c.Keyring().Key(c.L.ChainClient.Config.Key)
	if err != nil {
		panic(err)
	}
	address, err := key.GetAddress()
	if err != nil {
		panic(err)
	}
	return address
}

func getCreator(c Client) string {
	return c.addressString(getMainAddress(c))
}

func getCreatorAsValoper(c Client) string {
	return c.addressString(sdk.ValAddress(getMainAddress(c).Bytes()))
}

func (c Client) addressString(val sdk.Address) string {
	defer c.L.SetSDKContext()()
	return val.String()
}
