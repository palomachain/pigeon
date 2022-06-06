package relayer

import (
	"context"
	"testing"

	"github.com/palomachain/sparrow/attest"
	"github.com/palomachain/sparrow/chain"
	chainmocks "github.com/palomachain/sparrow/chain/mocks"
	"github.com/palomachain/sparrow/config"
	"github.com/palomachain/sparrow/relayer/mocks"
	"github.com/stretchr/testify/require"
)

func TestProcessing(t *testing.T) {
	ctx := context.Background()
	for _, tt := range []struct {
		name   string
		setup  func(t *testing.T) *Relayer
		expErr error
	}{
		{
			name: "without any processor it does nothing",
			setup: func(t *testing.T) *Relayer {
				return New(
					config.Root{},
					mocks.NewPalomaClienter(t),
					attest.NewRegistry(),
					make(map[string]chain.Processor),
				)
			},
		},
		{
			name: "it relays messages",
			setup: func(t *testing.T) *Relayer {
				p := chainmocks.NewProcessor(t)
				p.On("SupportedQueues").Return([]string{"a"})
				pal := mocks.NewPalomaClienter(t)
				pal.On("QueryMessagesInQueue", ctx, "a").Return(
					[]chain.MessageWithSignatures{
						{}, {},
					},
					nil,
				)
				p.On(
					"ProcessMessages",
					ctx,
					"a",
					[]chain.MessageWithSignatures{
						{}, {},
					},
				).Return(nil)
				return New(
					config.Root{},
					pal,
					attest.NewRegistry(),
					map[string]chain.Processor{
						"test": p,
					},
				)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			relayer := tt.setup(t)
			require.NoError(t, relayer.init())
			err := relayer.Process(ctx)

			require.ErrorIs(t, err, tt.expErr)
		})
	}
}
