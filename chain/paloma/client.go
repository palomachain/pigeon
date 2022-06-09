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
	valset "github.com/palomachain/sparrow/types/paloma/x/valset/types"
	"github.com/palomachain/sparrow/util/slice"
)

//go:generate mockery --name=MessageSender
type MessageSender interface {
	SendMsg(ctx context.Context, msg sdk.Msg) (*sdk.TxResponse, error)
}

type Client struct {
	L            *chain.LensClient
	palomaConfig config.Paloma

	GRPCClient grpc.ClientConn

	MessageSender MessageSender
}

// QueryMessagesForSigning returns a list of messages from a given queueTypeName that
// need to be signed by the provided validator given the valAddress.
func (c Client) QueryMessagesForSigning(
	ctx context.Context,
	valAddress sdk.ValAddress,
	queueTypeName string,
) ([]chain.QueuedMessage, error) {
	return queryMessagesForSigning(ctx, c.GRPCClient, c.L.Codec.Marshaler, valAddress, queueTypeName)
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
				ValAddress: vs.ValAddress,
				Signature:  vs.Signature,
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
	ID            uint64
	QueueTypeName string
	Signature     []byte
	ExtraData     []byte
}

// BroadcastMessageSignatures takes a list of signatures that need to be sent over to the chain.
// It build the message and sends it over.
func (c Client) BroadcastMessageSignatures(ctx context.Context, signatures ...BroadcastMessageSignatureIn) error {
	return broadcastMessageSignatures(ctx, c.MessageSender, signatures...)
}

// QueryValidatorInfo returns info about the validator.
func (c Client) QueryValidatorInfo(ctx context.Context, valAddr sdk.ValAddress) ([]*valset.ExternalChainInfo, error) {
	qc := valset.NewQueryClient(c.GRPCClient)
	valInfoRes, err := qc.ValidatorInfo(ctx, &valset.QueryValidatorInfoRequest{
		ValAddr: valAddr.String(),
	})
	if err != nil {
		if strings.Contains(err.Error(), "item not found in store") {
			return nil, nil
		}
		return nil, err
	}

	return valInfoRes.ChainInfos, nil
}

// TODO: this is only temporary for easier testing
func (c Client) DeleteJob(ctx context.Context, queueTypeName string, id uint64) error {
	key, err := c.Keyring().Key(c.L.ChainClient.Config.Key)
	if err != nil {
		return err
	}
	unlock := c.L.SetSDKContext()
	addr := key.GetAddress().String()
	unlock()
	_, err = c.MessageSender.SendMsg(ctx, &consensus.MsgDeleteJob{
		Creator:       addr,
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

	msg := &valset.MsgAddExternalChainInfoForValidator{}

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
	signatures ...BroadcastMessageSignatureIn,
) error {
	if len(signatures) == 0 {
		return nil
	}
	var signedMessages []*consensus.MsgAddMessagesSignatures_MsgSignedMessage
	for _, sig := range signatures {
		signedMessages = append(signedMessages, &consensus.MsgAddMessagesSignatures_MsgSignedMessage{
			Id:            sig.ID,
			QueueTypeName: sig.QueueTypeName,
			Signature:     sig.Signature,
			ExtraData:     sig.ExtraData,
		})
	}
	msg := &consensus.MsgAddMessagesSignatures{
		SignedMessages: signedMessages,
	}
	_, err := ms.SendMsg(ctx, msg)
	return err
}

func (c Client) Keyring() keyring.Keyring {
	return c.L.Keybase
}
