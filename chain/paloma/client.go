package paloma

import (
	"context"
	"strings"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	"github.com/strangelove-ventures/lens/client/query"
	"github.com/vizualni/whoops"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/config"
	consensus "github.com/palomachain/pigeon/types/paloma/x/consensus/types"
	"github.com/palomachain/pigeon/types/paloma/x/evm/types"
	evm "github.com/palomachain/pigeon/types/paloma/x/evm/types"
	valset "github.com/palomachain/pigeon/types/paloma/x/valset/types"
	"github.com/palomachain/pigeon/util/slice"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

type ResultStatus = coretypes.ResultStatus

//go:generate mockery --name=MessageSender
type MessageSender interface {
	SendMsg(ctx context.Context, msg sdk.Msg) (*sdk.TxResponse, error)
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
		var ptr consensus.Message
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

// QueryMessagesInQueue returns all messages that are currently in the queue.
func (c Client) QueryMessagesInQueue(ctx context.Context, queueTypeName string) ([]chain.MessageWithSignatures, error) {
	return queryMessagesInQueue(
		ctx,
		queueTypeName,
		c.valAddr,
		c.GRPCClient,
		c.L.Codec.Marshaler,
	)
}

func queryMessagesInQueue(
	ctx context.Context,
	queueTypeName string,
	skipEvidence sdk.ValAddress,
	c grpc.ClientConn,
	anyunpacker codectypes.AnyUnpacker,
) ([]chain.MessageWithSignatures, error) {
	qc := consensus.NewQueryClient(c)
	msgs, err := qc.MessagesInQueue(ctx, &consensus.QueryMessagesInQueueRequest{
		QueueTypeName:                    queueTypeName,
		SkipEvidenceProvidedByValAddress: skipEvidence,
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
		var ptr consensus.Message
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

func (c Client) QueryGetEVMValsetByID(ctx context.Context, id uint64, chainReferenceID string) (*types.Valset, error) {
	qc := evm.NewQueryClient(c.GRPCClient)
	valsetRes, err := qc.GetValsetByID(ctx, &evm.QueryGetValsetByIDRequest{
		ValsetID:         id,
		ChainReferenceID: chainReferenceID,
	})
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
	})
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

	_, err := c.MessageSender.SendMsg(ctx, msg)
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

	_, err = c.MessageSender.SendMsg(ctx, msg)
	return err
}

func (c Client) SetPublicAccessData(ctx context.Context, queueTypeName string, messageID uint64, data []byte) error {
	msg := &consensus.MsgSetPublicAccessData{
		Creator:       c.creator,
		Data:          data,
		MessageID:     messageID,
		QueueTypeName: queueTypeName,
	}

	_, err := c.MessageSender.SendMsg(ctx, msg)
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

	_, err := c.MessageSender.SendMsg(ctx, msg)
	return err
}

func (c Client) Status(ctx context.Context) (*ResultStatus, error) {
	q := query.Query{Client: &c.L.ChainClient, Options: query.DefaultOptions()}
	return q.Status()
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
	_, err := ms.SendMsg(ctx, msg)
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
	return key.GetAddress()
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
