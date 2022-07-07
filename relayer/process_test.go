package relayer

import (
	"context"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/evm"
	chainmocks "github.com/palomachain/pigeon/chain/mocks"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/relayer/mocks"
	evmtypes "github.com/palomachain/pigeon/types/paloma/x/evm/types"
	"github.com/stretchr/testify/mock"
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
				pc := mocks.NewPalomaClienter(t)
				pc.On("QueryGetEVMChainInfos", mock.Anything, mock.Anything).Return(nil, nil)
				return New(
					config.Root{},
					pc,
					evm.NewFactory(pc),
				)
			},
		},
		{
			name: "it relays messages",
			setup: func(t *testing.T) *Relayer {
				keyringPass := "abcd"

				dir := t.TempDir()
				keyring := evm.OpenKeystore(dir)
				acc, err := keyring.NewAccount(keyringPass)
				require.NoError(t, err)

				p := chainmocks.NewProcessor(t)
				p.On("SupportedQueues").Return([]string{"a"})
				p.On(
					"ProcessMessages",
					ctx,
					"a",
					[]chain.MessageWithSignatures{
						{}, {},
					},
				).Return(nil).Maybe() // todo: remove maybe later

				pal := mocks.NewPalomaClienter(t)
				pal.On("QueryGetEVMChainInfos", mock.Anything, mock.Anything).Return([]*evmtypes.ChainInfo{
					{
						ChainReferenceID:     "main",
						ChainID:              5,
						SmartContractID:      "5",
						SmartContractAddr:    common.BytesToAddress([]byte("abcd")).Hex(),
						ReferenceBlockHeight: 5,
						ReferenceBlockHash:   "0x12",
					},
				}, nil)
				pal.On("QueryMessagesInQueue", ctx, mock.Anything).Return(
					[]chain.MessageWithSignatures{
						{}, {},
					},
					nil,
				)
				pal.On("QueryMessagesForSigning", ctx, "a").Return(
					[]chain.QueuedMessage{},
					nil,
				)
				os.Setenv("TEST_PASS", keyringPass)
				t.Cleanup(func() {
					os.Unsetenv("TEST_PASS")
				})

				factory := mocks.NewEvmFactorier(t)

				factory.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(p, nil)

				return New(
					config.Root{
						EVM: map[string]config.EVM{
							"main": {
								ChainClientConfig: config.ChainClientConfig{
									KeyringPassEnvName: "TEST_PASS",
									SigningKey:         acc.Address.Hex(),
									KeyringDirectory:   config.Filepath(dir),
								},
							},
						},
					},
					pal,
					factory,
				)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			relayer := tt.setup(t)
			require.NoError(t, relayer.init())

			processors, err := relayer.buildProcessors(ctx)
			require.NoError(t, err)

			err = relayer.Process(ctx, processors)
			require.ErrorIs(t, err, tt.expErr)
		})
	}
}
