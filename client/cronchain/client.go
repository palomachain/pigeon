package cronchain

import (
	"context"
	"errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gogo/protobuf/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	chain "github.com/volumefi/conductor/client"
	cronchain "github.com/volumefi/conductor/types/cronchain"
)

type Client struct {
	L *chain.LensClient
}

type QueuedMessage[T cronchain.Signable] struct {
	ID    uint64
	Nonce []byte
	Msg   T
}

type grpcInvoker interface {
	Invoke()
}

// QueryMessagesForSigning returns a list of messages from a given queueTypeName that
// need to be signed by the provided validator given the valAddress.
func QueryMessagesForSigning[T cronchain.Signable](
	ctx context.Context,
	c Client,
	valAddress string,
	queueTypeName string,
) ([]QueuedMessage[T], error) {
	return queryMessagesForSigning[T](ctx, c.L, c.L.Codec.Marshaler, valAddress, queueTypeName)
}

func queryMessagesForSigning[T cronchain.Signable](
	ctx context.Context,
	c grpc.ClientConn,
	anyunpacker codectypes.AnyUnpacker,
	valAddress string,
	queueTypeName string,
) ([]QueuedMessage[T], error) {
	qc := cronchain.NewQueryClient(c)
	msgs, err := qc.QueuedMessagesForSigning(ctx, &cronchain.QueryQueuedMessagesForSigningRequest{
		ValAddress:    valAddress,
		QueueTypeName: queueTypeName,
	})
	if err != nil {
		return nil, err
	}
	var res []QueuedMessage[T]
	for _, msg := range msgs.GetMessageToSign() {
		var m cronchain.Signable
		err := anyunpacker.UnpackAny(msg.GetMsg(), &m)
		if err != nil {
			return nil, err
		}
		msgT, ok := m.(T)
		if !ok {
			return nil, errors.New("onmg")
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

func (c Client) BroadcastMessageSignatures(ctx context.Context, signatures ...BroadcastMessageSignatureIn) error {
	var signedMessages []*cronchain.MsgAddMessagesSignatures_MsgSignedMessage
	for _, sig := range signatures {
		signedMessages = append(signedMessages, &cronchain.MsgAddMessagesSignatures_MsgSignedMessage{
			Id:            sig.ID,
			QueueTypeName: sig.QueueTypeName,
			Signature:     sig.Signature,
		})
	}
	info, _ := c.Keyring().Key(c.L.Config.Key)
	addr, _ := c.L.DecodeBech32AccAddr(info.GetAddress().String())
	msg := &cronchain.MsgAddMessagesSignatures{
		Creator:        addr.String(),
		SignedMessages: signedMessages,
	}
	_, err := c.L.SendMsg(ctx, msg)
	return err
}

func (c Client) Keyring() keyring.Keyring {
	return c.L.Keybase
}
