package paloma

import (
	"context"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/grpc"
	"github.com/vizualni/whoops"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	chain "github.com/palomachain/sparrow/client"
	"github.com/palomachain/sparrow/config"
	paloma "github.com/palomachain/sparrow/types/paloma"
)

type Client struct {
	L            *chain.LensClient
	palomaConfig config.Paloma

	GRPCClient grpc.ClientConn
}

type QueuedMessage[T paloma.Signable] struct {
	ID    uint64
	Nonce []byte
	Msg   T
}

// QueryMessagesForSigning returns a list of messages from a given queueTypeName that
// need to be signed by the provided validator given the valAddress.
func QueryMessagesForSigning[T paloma.Signable](
	ctx context.Context,
	c Client,
	valAddress string,
	queueTypeName string,
) ([]QueuedMessage[T], error) {
	return queryMessagesForSigning[T](ctx, c.GRPCClient, c.L.Codec.Marshaler, valAddress, queueTypeName)
}

func queryMessagesForSigning[T paloma.Signable](
	ctx context.Context,
	c grpc.ClientConn,
	anyunpacker codectypes.AnyUnpacker,
	valAddress string,
	queueTypeName string,
) ([]QueuedMessage[T], error) {
	qc := paloma.NewQueryClient(c)
	msgs, err := qc.QueuedMessagesForSigning(ctx, &paloma.QueryQueuedMessagesForSigningRequest{
		ValAddress:    valAddress,
		QueueTypeName: queueTypeName,
	})
	if err != nil {
		return nil, err
	}
	var res []QueuedMessage[T]
	for _, msg := range msgs.GetMessageToSign() {
		var m paloma.Signable
		err := anyunpacker.UnpackAny(msg.GetMsg(), &m)
		if err != nil {
			return nil, whoops.Wrap(err, ErrUnableToUnpackAny)
		}
		msgT, ok := m.(T)
		if !ok {
			var expectedType T
			return nil, ErrIncorrectTypeSavedInMessage.Format(expectedType, m)
		}
		res = append(res, QueuedMessage[T]{
			ID:    msg.GetId(),
			Nonce: msg.GetNonce(),
			Msg:   msgT,
		})
	}

	return res, nil
}

type BroadcastMessageSignatureIn struct {
	ID            uint64
	QueueTypeName string
	Signature     []byte
}

// BroadcastMessageSignatures takes a list of signatures that need to be sent over to the chain.
// It build the message and sends it over.
func (c Client) BroadcastMessageSignatures(ctx context.Context, signatures ...BroadcastMessageSignatureIn) error {
	return broadcastMessageSignatures(ctx, c.L, signatures...)
}

// QueryValidatorInfo returns info about the validator.
func (c Client) QueryValidatorInfo(ctx context.Context) (*paloma.Validator, error) {
	qc := paloma.NewQueryValsetClient(c.GRPCClient)
	valInfoRes, err := qc.ValidatorInfo(ctx, &paloma.QueryValidatorInfoRequest{
		ValAddr: "TODO CHANGE ME", // TODO: pass in key!!! this is name
	})
	if err != nil {
		return nil, err
	}

	return valInfoRes.Validator, nil
}

// RegisterValidator registers itself with the network and sends it's public key that they are using for
// signing messages.
func (c Client) RegisterValidator(ctx context.Context, pubKey, signedPubKey []byte) error {
	txsvc := paloma.NewValsetTxServiceClient(c.GRPCClient)

	_, err := txsvc.RegisterConductor(ctx, &paloma.MsgRegisterConductor{
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
	txsvc := paloma.NewValsetTxServiceClient(c.GRPCClient)

	msg := &paloma.MsgAddExternalChainInfoForValidator{}

	for _, ci := range chainInfos {
		msg.ChainInfos = append(msg.ChainInfos, &paloma.MsgAddExternalChainInfoForValidator_ChainInfo{
			ChainID: ci.ChainID,
			Address: ci.AccAddress,
		})
	}

	_, err := txsvc.AddExternalChainInfoForValidator(ctx, msg)
	return err
}

type msgSender interface {
	SendMsg(ctx context.Context, msg sdk.Msg) (*sdk.TxResponse, error)
}

func broadcastMessageSignatures(
	ctx context.Context,
	ms msgSender,
	signatures ...BroadcastMessageSignatureIn,
) error {
	if len(signatures) == 0 {
		return nil
	}
	var signedMessages []*paloma.MsgAddMessagesSignatures_MsgSignedMessage
	for _, sig := range signatures {
		signedMessages = append(signedMessages, &paloma.MsgAddMessagesSignatures_MsgSignedMessage{
			Id:            sig.ID,
			QueueTypeName: sig.QueueTypeName,
			Signature:     sig.Signature,
		})
	}
	msg := &paloma.MsgAddMessagesSignatures{
		SignedMessages: signedMessages,
	}
	_, err := ms.SendMsg(ctx, msg)
	return err
}

func (c Client) Keyring() keyring.Keyring {
	return c.L.Keybase
}
