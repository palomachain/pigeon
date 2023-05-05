package paloma

import (
	"context"
	"github.com/VolumeFi/whoops"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/grpc"
	log "github.com/sirupsen/logrus"
	"time"

	ggrpc "google.golang.org/grpc"
)

var _ grpc.ClientConn = GRPCClientDowner{}

var _ MessageSender = MessageSenderDowner{}

type GRPCClientDowner struct {
	W grpc.ClientConn
}

type MessageSenderDowner struct {
	W MessageSender
}

func (g GRPCClientDowner) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...ggrpc.CallOption) error {
	err := g.W.Invoke(ctx, method, args, reply, opts...)
	if IsPalomaDown(err) {
		return whoops.Wrap(ErrPalomaIsDown, err)
	}
	return err
}

func (g GRPCClientDowner) NewStream(ctx context.Context, desc *ggrpc.StreamDesc, method string, opts ...ggrpc.CallOption) (ggrpc.ClientStream, error) {
	stream, err := g.W.NewStream(ctx, desc, method, opts...)

	if IsPalomaDown(err) {
		return nil, whoops.Wrap(ErrPalomaIsDown, err)
	}

	return stream, err
}

func (m MessageSenderDowner) SendMsg(ctx context.Context, msg sdk.Msg, memo string) (*sdk.TxResponse, error) {
	time.Sleep(10 * time.Second)
	log.Debug("Sending Msg: ", msg)
	res, err := m.W.SendMsg(ctx, msg, memo)

	if IsPalomaDown(err) {
		return nil, whoops.Wrap(ErrPalomaIsDown, err)
	}

	return res, err
}
