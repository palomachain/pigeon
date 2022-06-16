package paloma

import (
	"context"
	"strings"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/palomachain/sparrow/chain"
	"github.com/palomachain/sparrow/config"
	consensus "github.com/palomachain/sparrow/types/paloma/x/consensus/types"
	"github.com/palomachain/sparrow/types/paloma/x/evm/types"
	evm "github.com/palomachain/sparrow/types/paloma/x/evm/types"
	valset "github.com/palomachain/sparrow/types/paloma/x/valset/types"
	"github.com/palomachain/sparrow/util/slice"
)

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
		c.GRPCClient,
		c.L.Codec.Marshaler,
	)
}

func queryMessagesInQueue(
	ctx context.Context,
	queueTypeName string,
	c grpc.ClientConn,
	anyunpacker codectypes.AnyUnpacker,
) ([]chain.MessageWithSignatures, error) {
	qc := consensus.NewQueryClient(c)
	msgs, err := qc.MessagesInQueue(ctx, &consensus.QueryMessagesInQueueRequest{
		QueueTypeName: queueTypeName,
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
				ID:    msg.Id,
				Nonce: msg.Nonce,
				Msg:   ptr,
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
	ExtraData       []byte
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

func (c Client) QueryGetSnapshotByID(ctx context.Context, id uint64) (*valset.Snapshot, error) {
	qc := valset.NewQueryClient(c.GRPCClient)
	snapshotRes, err := qc.GetSnapshotByID(ctx, &valset.QueryGetSnapshotByIDRequest{
		SnapshotId: id,
	})
	if err != nil {
		if strings.Contains(err.Error(), "item not found in store") {
			return nil, nil
		}
		return nil, err
	}

	return snapshotRes.Snapshot, nil
}

func (c Client) QueryGetEVMValsetByID(ctx context.Context, id uint64, chainID string) (*types.Valset, error) {
	qc := evm.NewQueryClient(c.GRPCClient)
	valsetRes, err := qc.GetValsetByID(ctx, &evm.QueryGetValsetByIDRequest{
		ValsetID: id,
		ChainID:  chainID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "item not found in store") {
			return nil, nil
		}
		return nil, err
	}

	return valsetRes.Valset, nil
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
	ChainType  string
	ChainID    string
	AccAddress string
	PubKey     []byte
}

// AddExternalChainInfo adds info about the external chain. It adds the chain's
// account addresses that the sparrow knows about.
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
				ChainType: in.ChainType,
				ChainID:   in.ChainID,
				Address:   in.AccAddress,
				Pubkey:    in.PubKey,
			}
		},
	)

	_, err := c.MessageSender.SendMsg(ctx, msg)
	return err
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
	var signedMessages []*consensus.MsgAddMessagesSignatures_MsgSignedMessage
	for _, sig := range signatures {
		signedMessages = append(signedMessages, &consensus.MsgAddMessagesSignatures_MsgSignedMessage{
			Id:              sig.ID,
			QueueTypeName:   sig.QueueTypeName,
			Signature:       sig.Signature,
			ExtraData:       sig.ExtraData,
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
