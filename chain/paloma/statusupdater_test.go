package paloma

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/palomachain/paloma/x/evm/types"
	palomatypes "github.com/palomachain/paloma/x/paloma/types"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/paloma/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStatusUpdater(t *testing.T) {
	testmsg := &chain.MessageWithSignatures{
		QueuedMessage: chain.QueuedMessage{
			ID:               17,
			Nonce:            []byte("nonce"),
			BytesToSign:      []byte("bts"),
			PublicAccessData: []byte("pad"),
			ErrorData:        []byte("ed"),
			Msg: &evmtypes.Message{
				Action: &evmtypes.Message_SubmitLogicCall{
					SubmitLogicCall: &evmtypes.SubmitLogicCall{
						HexContractAddress: "0xABC",
						Abi:                []byte("abi"),
						Payload:            []byte("payload"),
						Deadline:           123,
					},
				},
			},
		},
		Signatures: []chain.ValidatorSignature{},
	}
	t.Run("with error during sending", func(t *testing.T) {
		m := mocks.NewMessageSender(t)
		client := Client{MessageSender: m}
		m.On("SendMsg", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("fail"))
		err := client.NewStatus().Debug(context.Background())
		require.Error(t, err)
	})

	for _, tt := range []struct {
		setup func() *palomatypes.MsgAddStatusUpdate
		name  string
		exec  func(c *Client) StatusUpdater
	}{
		{
			name: "minimal setup",
			setup: func() *palomatypes.MsgAddStatusUpdate {
				return &palomatypes.MsgAddStatusUpdate{
					Status: "",
					Args:   []palomatypes.MsgAddStatusUpdate_KeyValuePair{},
				}
			},
			exec: func(c *Client) StatusUpdater { return c.NewStatus() },
		},
		{
			name: "with chain reference ID",
			setup: func() *palomatypes.MsgAddStatusUpdate {
				return &palomatypes.MsgAddStatusUpdate{
					Status: "",
					Args: []palomatypes.MsgAddStatusUpdate_KeyValuePair{
						{
							Key:   "chain-reference-id",
							Value: "testchain",
						},
					},
				}
			},
			exec: func(c *Client) StatusUpdater { return c.NewStatus().WithChainReferenceID("testchain") },
		},
		{
			name: "with queue name",
			setup: func() *palomatypes.MsgAddStatusUpdate {
				return &palomatypes.MsgAddStatusUpdate{
					Status: "",
					Args: []palomatypes.MsgAddStatusUpdate_KeyValuePair{
						{
							Key:   "queue-type",
							Value: "testqueue",
						},
					},
				}
			},
			exec: func(c *Client) StatusUpdater { return c.NewStatus().WithQueueType("testqueue") },
		},
		{
			name: "with log",
			setup: func() *palomatypes.MsgAddStatusUpdate {
				return &palomatypes.MsgAddStatusUpdate{
					Status: "foobar",
					Args:   []palomatypes.MsgAddStatusUpdate_KeyValuePair{},
				}
			},
			exec: func(c *Client) StatusUpdater { return c.NewStatus().WithLog("foobar") },
		},
		{
			name: "with msg",
			setup: func() *palomatypes.MsgAddStatusUpdate {
				return &palomatypes.MsgAddStatusUpdate{
					Status: "",
					Args: []palomatypes.MsgAddStatusUpdate_KeyValuePair{
						{
							Key:   "message-type",
							Value: "*types.Message_SubmitLogicCall",
						}, {
							Key:   "message-id",
							Value: "17",
						}, {
							Key:   "message-nonce",
							Value: hexutil.Encode([]byte("nonce")),
						}, {
							Key:   "message-error-data",
							Value: hexutil.Encode([]byte("ed")),
						}, {
							Key:   "message-public-access-data",
							Value: hexutil.Encode([]byte("pad")),
						},
					},
				}
			},
			exec: func(c *Client) StatusUpdater {
				return c.NewStatus().WithMsg(testmsg)
			},
		},
		{
			name: "with all",
			setup: func() *palomatypes.MsgAddStatusUpdate {
				return &palomatypes.MsgAddStatusUpdate{
					Status: "foobar",
					Args: []palomatypes.MsgAddStatusUpdate_KeyValuePair{
						{
							Key:   "foo",
							Value: "bar",
						},
						{
							Key:   "chain-reference-id",
							Value: "testchain",
						},
						{
							Key:   "queue-type",
							Value: "testqueue",
						},
						{
							Key:   "message-type",
							Value: "*types.Message_SubmitLogicCall",
						},
						{
							Key:   "message-id",
							Value: "17",
						},
						{
							Key:   "message-nonce",
							Value: hexutil.Encode([]byte("nonce")),
						},
						{
							Key:   "message-error-data",
							Value: hexutil.Encode([]byte("ed")),
						},
						{
							Key:   "message-public-access-data",
							Value: hexutil.Encode([]byte("pad")),
						},
					},
				}
			},
			exec: func(c *Client) StatusUpdater {
				return c.NewStatus().WithMsg(testmsg).WithArg("foo", "bar").WithChainReferenceID("testchain").WithQueueType("testqueue").WithLog("foobar")
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, cfg := range []struct {
				exec  func(StatusUpdater)
				name  string
				level palomatypes.MsgAddStatusUpdate_Level
			}{
				{
					name:  "on DEBUG",
					level: palomatypes.MsgAddStatusUpdate_LEVEL_DEBUG,
					exec:  func(s StatusUpdater) { s.Debug(context.Background()) },
				},
				{
					name:  "on INFO",
					level: palomatypes.MsgAddStatusUpdate_LEVEL_INFO,
					exec:  func(s StatusUpdater) { s.Info(context.Background()) },
				},
				{
					name:  "on ERROR",
					level: palomatypes.MsgAddStatusUpdate_LEVEL_ERROR,
					exec:  func(s StatusUpdater) { s.Error(context.Background()) },
				},
			} {
				t.Run(cfg.name, func(t *testing.T) {
					ms := mocks.NewMessageSender(t)
					msg := tt.setup()
					msg.Level = cfg.level
					ms.On("SendMsg", mock.Anything, msg, mock.Anything).Return(nil, nil)
					client := Client{MessageSender: ms}
					cfg.exec(tt.exec(&client))
				})
			}
		})
	}
}
