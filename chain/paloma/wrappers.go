package paloma

import (
	"context"

	"github.com/VolumeFi/whoops"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/grpc"
	"github.com/palomachain/pigeon/internal/liblog"
	ggrpc "google.golang.org/grpc"
)

var _ grpc.ClientConn = GRPCClientWrapper{}

var _ MessageSender = PalomaMessageSender{}

type GRPCClientWrapper struct {
	W grpc.ClientConn
}

type KeyRotator interface {
	RotateKeys(context.Context)
}

type PalomaMessageSender struct {
	R KeyRotator
	W MessageSender
}

func (g GRPCClientWrapper) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...ggrpc.CallOption) error {
	err := g.W.Invoke(ctx, method, args, reply, opts...)
	if IsPalomaDown(err) {
		return whoops.Wrap(ErrPalomaIsDown, err)
	}
	return err
}

func (g GRPCClientWrapper) NewStream(ctx context.Context, desc *ggrpc.StreamDesc, method string, opts ...ggrpc.CallOption) (ggrpc.ClientStream, error) {
	stream, err := g.W.NewStream(ctx, desc, method, opts...)

	if IsPalomaDown(err) {
		return nil, whoops.Wrap(ErrPalomaIsDown, err)
	}

	return stream, err
}

func (m PalomaMessageSender) SendMsg(ctx context.Context, msg sdk.Msg, memo string) (*sdk.TxResponse, error) {
	logger := liblog.WithContext(ctx).WithField("component", "message-sender").WithField("msg", msg)
	logger.Debug("Sending Msg")

	// TODO: use lock
	m.R.RotateKeys(ctx)

	res, err := m.W.SendMsg(ctx, msg, memo)
	if IsPalomaDown(err) {
		return nil, whoops.Wrap(ErrPalomaIsDown, err)
	}

	return res, err
}
