package paloma

import (
	"context"
	"errors"
	"net"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	clientmocks "github.com/palomachain/sparrow/client/paloma/mocks"
	consensus "github.com/palomachain/sparrow/types/paloma/x/consensus/types"
	consensusmocks "github.com/palomachain/sparrow/types/paloma/x/consensus/types/mocks"
	valset "github.com/palomachain/sparrow/types/paloma/x/valset/types"
	valsetmocks "github.com/palomachain/sparrow/types/paloma/x/valset/types/mocks"
	"github.com/palomachain/sparrow/types/testdata"
	"github.com/strangelove-ventures/lens/byop"
	lens "github.com/strangelove-ventures/lens/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

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

func queryServerDailer(t *testing.T, msgsrv *consensusmocks.QueryServer) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	consensus.RegisterQueryServer(server, msgsrv)

	go func() {
		err := server.Serve(listener)
		assert.NoError(t, err)
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func valsetQueryServerDailer(t *testing.T, msgsrv *valsetmocks.QueryServer) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	valset.RegisterQueryServer(server, msgsrv)

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
					Iface: (*consensus.Signable)(nil),
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
		name   string
		mcksrv func(*testing.T) *consensusmocks.QueryServer
		expRes []QueuedMessage[*testdata.SimpleMessage]
		expErr error

		// used only for testing the GRPC responses because GRPC is doing a
		// string concatenation on errors, thus we can't do proper error
		// inspection
		expectsAnyError bool
	}{
		{
			name: "called with correct arguments",
			mcksrv: func(t *testing.T) *consensusmocks.QueryServer {
				srv := consensusmocks.NewQueryServer(t)
				srv.On("QueuedMessagesForSigning", mock.Anything, &consensus.QueryQueuedMessagesForSigningRequest{
					ValAddress:    "validator",
					QueueTypeName: "queueName",
				}).Return(
					&consensus.QueryQueuedMessagesForSigningResponse{
						MessageToSign: []*consensus.MessageToSign{},
					},
					nil,
				).Once()
				return srv
			},
		},
		{
			name: "messages are happily returned",
			mcksrv: func(t *testing.T) *consensusmocks.QueryServer {
				msgany1, err := codectypes.NewAnyWithValue(simpleMessageTestData1)
				assert.NoError(t, err)
				msgany2, err := codectypes.NewAnyWithValue(simpleMessageTestData2)
				assert.NoError(t, err)

				srv := consensusmocks.NewQueryServer(t)
				srv.On("QueuedMessagesForSigning", mock.Anything, mock.Anything).Return(
					&consensus.QueryQueuedMessagesForSigningResponse{
						MessageToSign: []*consensus.MessageToSign{
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
					},
					nil,
				).Once()
				return srv
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
			mcksrv: func(t *testing.T) *consensusmocks.QueryServer {
				erroneousMsg := &codectypes.Any{
					TypeUrl: "/wrong",
					Value:   []byte(`whoops`),
				}

				srv := consensusmocks.NewQueryServer(t)
				srv.On("QueuedMessagesForSigning", mock.Anything, mock.Anything).Return(
					&consensus.QueryQueuedMessagesForSigningResponse{
						MessageToSign: []*consensus.MessageToSign{
							{
								Nonce: []byte("nonce-123"),
								Id:    456,
								Msg:   erroneousMsg,
							},
						},
					},
					nil,
				).Once()
				return srv
			},
			expErr: ErrUnableToUnpackAny,
		},
		{
			name: "client returns an error",
			mcksrv: func(t *testing.T) *consensusmocks.QueryServer {
				srv := consensusmocks.NewQueryServer(t)
				srv.On("QueuedMessagesForSigning", mock.Anything, mock.Anything).Return(
					nil,
					errTestErr,
				).Once()
				return srv
			},
			expectsAnyError: true,
		},
		{
			name: "incorrect type used for the unpacked message",
			mcksrv: func(t *testing.T) *consensusmocks.QueryServer {
				msgany, err := codectypes.NewAnyWithValue(&testdata.SimpleMessage2{
					Field: "random value",
				})
				assert.NoError(t, err)

				srv := consensusmocks.NewQueryServer(t)
				srv.On("QueuedMessagesForSigning", mock.Anything, mock.Anything).Return(
					&consensus.QueryQueuedMessagesForSigningResponse{
						MessageToSign: []*consensus.MessageToSign{
							{
								Nonce: []byte("nonce-123"),
								Id:    456,
								Msg:   msgany,
							},
						},
					},
					nil,
				).Once()
				return srv
			},
			expErr: ErrIncorrectTypeSavedInMessage,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// setting everything up
			ctx := context.Background()
			mocksrv := tt.mcksrv(t)
			conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(queryServerDailer(t, mocksrv)))
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

func TestRegisterValidator(t *testing.T) {
	pk := []byte{1, 2, 3}
	sig := []byte{4, 5, 6}
	fakeErr := errors.New("something")
	for _, tt := range []struct {
		name   string
		mcksrv func(*testing.T) *clientmocks.MessageSender
		expRes []QueuedMessage[*testdata.SimpleMessage]

		expectsAnyError bool
	}{
		{
			name: "happy path",
			mcksrv: func(t *testing.T) *clientmocks.MessageSender {
				srv := clientmocks.NewMessageSender(t)
				srv.On("SendMsg", mock.Anything, &valset.MsgRegisterConductor{
					Creator:      "signer1",
					ValAddr:      "val1",
					PubKey:       pk,
					SignedPubKey: sig,
				}).Return(nil, nil).Once()
				return srv
			},
		},
		{
			name: "grpc returns error",
			mcksrv: func(t *testing.T) *clientmocks.MessageSender {
				srv := clientmocks.NewMessageSender(t)
				srv.On("SendMsg", mock.Anything, mock.Anything).Return(nil, fakeErr).Once()
				return srv
			},
			expectsAnyError: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// setting everything up
			ctx := context.Background()
			mocksrv := tt.mcksrv(t)

			client := Client{
				MessageSender: mocksrv,
			}
			err := client.RegisterValidator(ctx, "signer1", "val1", pk, sig)

			if tt.expectsAnyError {
				assert.Error(t, err)
			}
		})
	}
}

func TestQueryValidatorInfo(t *testing.T) {
	fakeErr := errors.New("something")
	fakeVal := &valset.Validator{
		Address: "hello",
	}
	for _, tt := range []struct {
		name   string
		mcksrv func(*testing.T) *valsetmocks.QueryServer
		expRes []QueuedMessage[*testdata.SimpleMessage]

		expectedValidator *valset.Validator
		expectsAnyError   bool
	}{
		{
			name: "happy path",
			mcksrv: func(t *testing.T) *valsetmocks.QueryServer {
				srv := valsetmocks.NewQueryServer(t)
				srv.On("ValidatorInfo", mock.Anything, mock.Anything).Return(&valset.QueryValidatorInfoResponse{
					Validator: fakeVal,
				}, nil).Once()
				return srv
			},
			expectedValidator: fakeVal,
		},
		{
			name: "grpc returns error",
			mcksrv: func(t *testing.T) *valsetmocks.QueryServer {
				srv := valsetmocks.NewQueryServer(t)
				srv.On("ValidatorInfo", mock.Anything, mock.Anything).Return(nil, fakeErr).Once()
				return srv
			},
			expectsAnyError: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// setting everything up
			ctx := context.Background()
			mocksrv := tt.mcksrv(t)
			conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(valsetQueryServerDailer(t, mocksrv)))
			assert.NoError(t, err)

			client := Client{
				GRPCClient: conn,
			}
			valInfo, err := client.QueryValidatorInfo(ctx, "something")

			require.Equal(t, tt.expectedValidator, valInfo)

			if tt.expectsAnyError {
				assert.Error(t, err)
			}
		})
	}
}

func TestAddingExternalChainInfo(t *testing.T) {
	fakeErr := errors.New("something")
	for _, tt := range []struct {
		name      string
		chainInfo []ChainInfoIn
		mcksrv    func(*testing.T) *clientmocks.MessageSender
		expRes    []QueuedMessage[*testdata.SimpleMessage]

		expectsAnyError bool
	}{
		{
			name:      "without chain infos provided does nothing",
			chainInfo: []ChainInfoIn{},
			mcksrv: func(t *testing.T) *clientmocks.MessageSender {
				srv := clientmocks.NewMessageSender(t)
				t.Cleanup(func() {
					srv.AssertNotCalled(t, "SendMsg", mock.Anything, mock.Anything)
				})
				return srv
			},
		},
		{
			name: "happy path",
			chainInfo: []ChainInfoIn{
				{ChainID: "chain1", AccAddress: "addr1"},
				{ChainID: "chain2", AccAddress: "addr2"},
			},
			mcksrv: func(t *testing.T) *clientmocks.MessageSender {
				srv := clientmocks.NewMessageSender(t)
				srv.On("SendMsg", mock.Anything, &valset.MsgAddExternalChainInfoForValidator{
					ChainInfos: []*valset.MsgAddExternalChainInfoForValidator_ChainInfo{
						{ChainID: "chain1", Address: "addr1"},
						{ChainID: "chain2", Address: "addr2"},
					},
				}).Return(nil, nil).Once()
				return srv
			},
		},
		{
			name: "with SendMsg returning errors",
			chainInfo: []ChainInfoIn{
				{ChainID: "chain1", AccAddress: "addr1"},
				{ChainID: "chain2", AccAddress: "addr2"},
			},
			mcksrv: func(t *testing.T) *clientmocks.MessageSender {
				srv := clientmocks.NewMessageSender(t)
				srv.On("SendMsg", mock.Anything, mock.Anything).Return(nil, fakeErr).Once()
				return srv
			},
			expectsAnyError: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// setting everything up
			ctx := context.Background()
			mocksrv := tt.mcksrv(t)

			client := Client{
				MessageSender: mocksrv,
			}
			err := client.AddExternalChainInfo(
				ctx,
				tt.chainInfo...,
			)

			if tt.expectsAnyError {
				assert.Error(t, err)
			}
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
		msgSender  MessageSender
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
				addMsgSigs, ok := msg.(*consensus.MsgAddMessagesSignatures)
				assert.True(t, ok, "incorrect msg type")
				assert.Len(t, addMsgSigs.SignedMessages, 2)
				assert.Equal(t, addMsgSigs.SignedMessages[0], &consensus.MsgAddMessagesSignatures_MsgSignedMessage{
					Id:            123,
					QueueTypeName: "abc",
					Signature:     []byte(`sig-123`),
				})
				assert.Equal(t, addMsgSigs.SignedMessages[1], &consensus.MsgAddMessagesSignatures_MsgSignedMessage{
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
