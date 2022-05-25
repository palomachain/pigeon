package paloma

import (
	"context"
	"fmt"
	"strings"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/palomachain/sparrow/chain"
	"github.com/palomachain/sparrow/config"
	consensus "github.com/palomachain/sparrow/types/paloma/x/consensus/types"
	valset "github.com/palomachain/sparrow/types/paloma/x/valset/types"
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
		var ptr consensus.Signable
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

func QueryConsensusReachedMessages(
	ctx context.Context,
	c Client,
	queueTypeName string,
) ([]chain.ConsensusReachedMsg, error) {
	return queryConsensusReachedMessages(ctx, c.GRPCClient, c.L.Codec.Marshaler, queueTypeName)
}

func queryConsensusReachedMessages(
	ctx context.Context,
	c grpc.ClientConn,
	anyunpacker codectypes.AnyUnpacker,
	queueTypeName string,
) ([]chain.ConsensusReachedMsg, error) {
	qc := consensus.NewQueryClient(c)
	consensusRes, err := qc.ConsensusReached(ctx, &consensus.QueryConsensusReachedRequest{
		QueueTypeName: queueTypeName,
	})

	if err != nil {
		return nil, err
	}

	var res []chain.ConsensusReachedMsg
	for _, rawMsg := range consensusRes.GetMessages() {
		var signable consensus.Signable
		err := anyunpacker.UnpackAny(rawMsg.GetMsg(), &signable)
		if err != nil {
			return nil, err
		}

		m := chain.ConsensusReachedMsg{
			ID:    fmt.Sprintf("%d", rawMsg.GetId()),
			Nonce: rawMsg.Nonce,
			Msg:   signable,
		}
		for _, signData := range rawMsg.GetSignData() {
			m.Signatures = append(m.Signatures, chain.ValidatorSignature{
				ValAddress: signData.GetValAddress(),
				Signature:  signData.GetSignature(),
			})
		}
		res = append(res, m)
	}

	return res, nil
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
func (c Client) QueryValidatorInfo(ctx context.Context, valAddr string) (*valset.Validator, error) {
	qc := valset.NewQueryClient(c.GRPCClient)
	valInfoRes, err := qc.ValidatorInfo(ctx, &valset.QueryValidatorInfoRequest{
		ValAddr: valAddr,
	})
	if err != nil {
		if strings.Contains(err.Error(), "item not found in store") {
			return nil, nil
		}
		return nil, err
	}

	return valInfoRes.Validator, nil
}

// RegisterValidator registers itself with the network and sends it's public key that they are using for
// signing messages.
func (c Client) RegisterValidator(ctx context.Context, signerAddr, valAddr string, pubKey, signedPubKey []byte) error {
	_, err := c.MessageSender.SendMsg(ctx, &valset.MsgRegisterConductor{
		Creator:      signerAddr,
		ValAddr:      valAddr,
		PubKey:       pubKey,
		SignedPubKey: signedPubKey,
	})

	return err
}

type ChainInfoIn struct {
	ChainID    string
	AccAddress string
}

// AddExternalChainInfo adds info about the external chain. It adds the chain's
// account addresses that the runner owns.
func (c Client) AddExternalChainInfo(ctx context.Context, chainInfos ...ChainInfoIn) error {
	if len(chainInfos) == 0 {
		return nil
	}

	msg := &valset.MsgAddExternalChainInfoForValidator{}

	for _, ci := range chainInfos {
		msg.ChainInfos = append(msg.ChainInfos, &valset.MsgAddExternalChainInfoForValidator_ChainInfo{
			ChainID: ci.ChainID,
			Address: ci.AccAddress,
		})
	}

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
