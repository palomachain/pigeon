package paloma

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/palomachain/paloma/x/evm/types"
	palomatypes "github.com/palomachain/paloma/x/paloma/types"
	"github.com/palomachain/pigeon/chain"
)

var _ StatusUpdater = (*statusUpdater)(nil)

type StatusUpdater interface {
	WithLog(status string) StatusUpdater
	WithMsg(msg *chain.MessageWithSignatures) StatusUpdater
	WithQueueType(queueType string) StatusUpdater
	WithChainReferenceID(chainReferenceID string) StatusUpdater
	WithArg(key, value string) StatusUpdater
	Info(ctx context.Context) error
	Error(ctx context.Context) error
	Debug(ctx context.Context) error
}

func (c *Client) NewStatus() StatusUpdater {
	return &statusUpdater{
		c:    c,
		args: make(map[string]string),
	}
}

type statusUpdater struct {
	msg    *chain.MessageWithSignatures
	c      *Client
	args   map[string]string
	status string
	level  palomatypes.MsgAddStatusUpdate_Level
}

func (s *statusUpdater) WithLog(status string) StatusUpdater {
	s.status = status
	return s
}

func (s *statusUpdater) WithMsg(msg *chain.MessageWithSignatures) StatusUpdater {
	s.msg = msg
	return s
}

func (s *statusUpdater) WithQueueType(queueType string) StatusUpdater {
	s.args["queue-type"] = queueType
	return s
}

func (s *statusUpdater) WithChainReferenceID(chainReferenceID string) StatusUpdater {
	s.args["chain-reference-id"] = chainReferenceID
	return s
}

func (s *statusUpdater) WithArg(key, value string) StatusUpdater {
	s.args[key] = value
	return s
}

func (s *statusUpdater) Info(ctx context.Context) error {
	s.level = palomatypes.MsgAddStatusUpdate_LEVEL_INFO
	return s.update(ctx)
}

func (s *statusUpdater) Error(ctx context.Context) error {
	s.level = palomatypes.MsgAddStatusUpdate_LEVEL_ERROR
	return s.update(ctx)
}

func (s *statusUpdater) Debug(ctx context.Context) error {
	s.level = palomatypes.MsgAddStatusUpdate_LEVEL_DEBUG
	return s.update(ctx)
}

func (s *statusUpdater) update(ctx context.Context) error {
	args := make([]palomatypes.MsgAddStatusUpdate_KeyValuePair, 0, len(s.args)+3)
	for k, v := range s.args {
		args = append(args, palomatypes.MsgAddStatusUpdate_KeyValuePair{Key: k, Value: v})
	}

	if s.msg != nil {
		m := s.msg.Msg.(*evmtypes.Message)
		args = append(args,
			palomatypes.MsgAddStatusUpdate_KeyValuePair{
				Key:   "message-type",
				Value: fmt.Sprintf("%T", m.GetAction()),
			}, palomatypes.MsgAddStatusUpdate_KeyValuePair{
				Key:   "message-id",
				Value: fmt.Sprintf("%d", s.msg.ID),
			}, palomatypes.MsgAddStatusUpdate_KeyValuePair{
				Key:   "message-nonce",
				Value: hexutil.Encode(s.msg.Nonce),
			}, palomatypes.MsgAddStatusUpdate_KeyValuePair{
				Key:   "message-error-data",
				Value: hexutil.Encode(s.msg.ErrorData),
			}, palomatypes.MsgAddStatusUpdate_KeyValuePair{
				Key:   "message-public-access-data",
				Value: hexutil.Encode(s.msg.PublicAccessData),
			})
	}

	msg := &palomatypes.MsgAddStatusUpdate{
		Status: s.status,
		Level:  s.level,
		Args:   args,
	}

	_, err := s.c.MessageSender.SendMsg(ctx, msg, "", s.c.sendingOpts...)
	return err
}
