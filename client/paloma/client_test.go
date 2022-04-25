package paloma

import (
	"context"
	"errors"
	"net"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/strangelove-ventures/lens/byop"
	lens "github.com/strangelove-ventures/lens/client"
	"github.com/stretchr/testify/assert"
	"github.com/volumefi/conductor/types/cronchain"
	"github.com/volumefi/conductor/types/testdata"

	"google.golang.org/grpc/test/bufconn"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"google.golang.org/grpc"
)

var (
	errTestErr             = errors.New("sample error")
	simpleMessageTestData1 = &testdata.SimpleMessage{
		Sender: "bob",
		Hello:  "mars",
		World:  "!",
	}
	simpleMessageTestData2 = &testdata.SimpleMessage{
		Sender: "alice",
		Hello:  "jupiter",
		World:  "!",
	}
)

var _ cronchain.QueryServer = mockCronchainQueryServer{}

type queuedMsgsFnc func(context.Context, *cronchain.QueryQueuedMessagesForSigningRequest) (*cronchain.QueryQueuedMessagesForSigningResponse, error)

type mockCronchainQueryServer struct {
	queuedMsgs queuedMsgsFnc
}

func (m mockCronchainQueryServer) QueuedMessagesForSigning(ctx context.Context, msg *cronchain.QueryQueuedMessagesForSigningRequest) (*cronchain.QueryQueuedMessagesForSigningResponse, error) {
	return m.queuedMsgs(ctx, msg)
}

func dialer(t *testing.T, msgsrv mockCronchainQueryServer) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	cronchain.RegisterQueryServer(server, msgsrv)

	go func() {
		err := server.Serve(listener)
		assert.NoError(t, err)
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func makeCodec() lens.Codec {
	return lens.MakeCodec([]module.AppModuleBasic{
		byop.Module{
			ModuleName: "testing",
			MsgsImplementations: []byop.RegisterImplementation{
				{
					Iface: (*cronchain.Signable)(nil),
					Msgs: []proto.Message{
						&testdata.SimpleMessage{},
						&testdata.SimpleMessage2{},
					},
				},
			},
		},
	})
}

func TestQueryingMessagesForSigning(t *testing.T) {
	codec := makeCodec()
	for _, tt := range []struct {
		name          string
		queuedMsgsFnc queuedMsgsFnc
		expRes        []QueuedMessage[*testdata.SimpleMessage]
		expErr        error

		// used only for testing the GRPC responses because GRPC is doing a
		// string concatenation on errors, thus we can't do proper error
		// inspection
		expectsAnyError bool
	}{
		{
			name: "called with correct arguments",
			queuedMsgsFnc: func(_ context.Context, msg *cronchain.QueryQueuedMessagesForSigningRequest) (*cronchain.QueryQueuedMessagesForSigningResponse, error) {
				assert.Equal(t, "validator", msg.ValAddress)
				assert.Equal(t, "queueName", msg.QueueTypeName)
				return &cronchain.QueryQueuedMessagesForSigningResponse{}, nil
			},
		},
		{
			name: "messages are happily returned",
			queuedMsgsFnc: func(_ context.Context, msg *cronchain.QueryQueuedMessagesForSigningRequest) (*cronchain.QueryQueuedMessagesForSigningResponse, error) {
				msgany1, err := codectypes.NewAnyWithValue(simpleMessageTestData1)
				assert.NoError(t, err)
				msgany2, err := codectypes.NewAnyWithValue(simpleMessageTestData2)
				assert.NoError(t, err)
				return &cronchain.QueryQueuedMessagesForSigningResponse{
					MessageToSign: []*cronchain.MessageToSign{
						{
							Nonce: []byte("nonce-123"),
							Id:    456,
							Msg:   msgany1,
						},
						{
							Nonce: []byte("nonce-321"),
							Id:    654,
							Msg:   msgany2,
						},
					},
				}, nil
			},
			expRes: []QueuedMessage[*testdata.SimpleMessage]{
				{
					Nonce: []byte("nonce-123"),
					ID:    456,
					Msg:   simpleMessageTestData1,
				},
				{
					Nonce: []byte("nonce-321"),
					ID:    654,
					Msg:   simpleMessageTestData2,
				},
			},
		},
		{
			name: "unpacking messages returns an error",
			queuedMsgsFnc: func(_ context.Context, msg *cronchain.QueryQueuedMessagesForSigningRequest) (*cronchain.QueryQueuedMessagesForSigningResponse, error) {
				erroneousMsg := &codectypes.Any{
					TypeUrl: "/wrong",
					Value:   []byte(`whoops`),
				}
				return &cronchain.QueryQueuedMessagesForSigningResponse{
					MessageToSign: []*cronchain.MessageToSign{
						{
							Nonce: []byte("nonce-123"),
							Id:    456,
							Msg:   erroneousMsg,
						},
					},
				}, nil
			},
			expErr: ErrUnableToUnpackAny,
		},
		{
			name: "client returns an error",
			queuedMsgsFnc: func(_ context.Context, msg *cronchain.QueryQueuedMessagesForSigningRequest) (*cronchain.QueryQueuedMessagesForSigningResponse, error) {
				return nil, errTestErr
			},
			expectsAnyError: true,
		},
		{
			name: "incorrect type used for the unpacked message",
			queuedMsgsFnc: func(_ context.Context, msg *cronchain.QueryQueuedMessagesForSigningRequest) (*cronchain.QueryQueuedMessagesForSigningResponse, error) {
				msgany, err := codectypes.NewAnyWithValue(&testdata.SimpleMessage2{
					Field: "random value",
				})
				assert.NoError(t, err)
				return &cronchain.QueryQueuedMessagesForSigningResponse{
					MessageToSign: []*cronchain.MessageToSign{
						{
							Nonce: []byte("nonce-123"),
							Id:    456,
							Msg:   msgany,
						},
					},
				}, nil
			},
			expErr: ErrIncorrectTypeSavedInMessage,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// setting everything up
			ctx := context.Background()
			mocksrv := mockCronchainQueryServer{}
			mocksrv.queuedMsgs = tt.queuedMsgsFnc
			conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(t, mocksrv)))
			assert.NoError(t, err)

			// setup complete
			// calling the function that we are testing
			msgs, err := queryMessagesForSigning[*testdata.SimpleMessage](
				ctx,
				conn,
				codec.Marshaler,
				"validator",
				"queueName",
			)
			if tt.expectsAnyError {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tt.expErr)
			}
			assert.Equal(t, tt.expRes, msgs)
		})
	}
}

type mockMsgSender func(context.Context, sdk.Msg) (*sdk.TxResponse, error)

func (m mockMsgSender) SendMsg(ctx context.Context, msg sdk.Msg) (*sdk.TxResponse, error) {
	return m(ctx, msg)
}
func TestBroadcastingMessageSignatures(t *testing.T) {
	ctx := context.Background()
	for _, tt := range []struct {
		name       string
		msgSender  msgSender
		signatures []BroadcastMessageSignatureIn

		expErr error
	}{
		{
			name: "nothing happens when there are no signatures being sent",
		},
		{
			name: "signatures are sent over",
			signatures: []BroadcastMessageSignatureIn{
				{
					ID:            123,
					QueueTypeName: "abc",
					Signature:     []byte(`sig-123`),
				},
				{
					ID:            456,
					QueueTypeName: "def",
					Signature:     []byte(`sig-789`),
				},
			},
			msgSender: mockMsgSender(func(ctx context.Context, msg sdk.Msg) (*sdk.TxResponse, error) {
				addMsgSigs, ok := msg.(*cronchain.MsgAddMessagesSignatures)
				assert.True(t, ok, "incorrect msg type")
				assert.Len(t, addMsgSigs.SignedMessages, 2)
				assert.Equal(t, addMsgSigs.SignedMessages[0], &cronchain.MsgAddMessagesSignatures_MsgSignedMessage{
					Id:            123,
					QueueTypeName: "abc",
					Signature:     []byte(`sig-123`),
				})
				assert.Equal(t, addMsgSigs.SignedMessages[1], &cronchain.MsgAddMessagesSignatures_MsgSignedMessage{
					Id:            456,
					QueueTypeName: "def",
					Signature:     []byte(`sig-789`),
				})
				return nil, nil
			}),
		},
		{
			name: "msg sender returns an error",
			msgSender: mockMsgSender(func(ctx context.Context, msg sdk.Msg) (*sdk.TxResponse, error) {
				return nil, errTestErr
			}),
			signatures: []BroadcastMessageSignatureIn{
				{},
			},
			expErr: errTestErr,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := broadcastMessageSignatures(
				ctx,
				tt.msgSender,
				tt.signatures...,
			)
			assert.ErrorIs(t, tt.expErr, err)
		})
	}
}
