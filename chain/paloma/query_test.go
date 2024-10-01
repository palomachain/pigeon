package paloma

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	consensus "github.com/palomachain/paloma/v2/x/consensus/types"
	consensusmocks "github.com/palomachain/paloma/v2/x/consensus/types/mocks"
	"github.com/palomachain/pigeon/chain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestQueryingMessagesForSigning(t *testing.T) {
	codec := makeCodec()
	for _, tt := range []struct {
		name   string
		mcksrv func(*testing.T) *consensusmocks.QueryServer
		expRes []chain.QueuedMessage
		expErr error

		// used only for testing the GRPC responses because GRPC is doing a
		// string concatenation on errors, thus we can't do proper error
		// inspection
		expectsAnyError bool
	}{
		{
			name:   "called with correct arguments",
			expRes: []chain.QueuedMessage{},
			mcksrv: func(t *testing.T) *consensusmocks.QueryServer {
				srv := consensusmocks.NewQueryServer(t)
				srv.On("QueuedMessagesForSigning", mock.Anything, &consensus.QueryQueuedMessagesForSigningRequest{
					ValAddress:    sdk.ValAddress("validator"),
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
				srv := consensusmocks.NewQueryServer(t)
				srv.On("QueuedMessagesForSigning", mock.Anything, mock.Anything).Return(
					&consensus.QueryQueuedMessagesForSigningResponse{
						MessageToSign: []*consensus.MessageToSign{
							{
								Nonce:       []byte("nonce-123"),
								Id:          456,
								BytesToSign: []byte("bla"),
							},
							{
								Nonce:       []byte("nonce-321"),
								Id:          654,
								BytesToSign: []byte("bla2"),
							},
						},
					},
					nil,
				).Once()
				return srv
			},
			expRes: []chain.QueuedMessage{
				{
					Nonce:       []byte("nonce-123"),
					ID:          456,
					BytesToSign: []byte("bla"),
				},
				{
					Nonce:       []byte("nonce-321"),
					ID:          654,
					BytesToSign: []byte("bla2"),
				},
			},
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
	} {
		t.Run(tt.name, func(t *testing.T) {
			// setting everything up
			ctx := context.Background()
			mocksrv := tt.mcksrv(t)
			conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(consensusQueryServerDialer(t, mocksrv)))
			require.NoError(t, err)

			// setup complete
			// calling the function that we are testing
			msgs, err := queryMessagesForSigning(
				ctx,
				conn,
				codec.Marshaler,
				sdk.ValAddress("validator"),
				"queueName",
			)
			if tt.expectsAnyError {
				require.Error(t, err)
			} else {
				require.ErrorIs(t, err, tt.expErr)
			}
			require.Equal(t, tt.expRes, msgs)
		})
	}
}
