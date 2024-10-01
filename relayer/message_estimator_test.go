package relayer

import (
	"context"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/palomachain/paloma/v2/x/evm/types"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/evm"
	evmmocks "github.com/palomachain/pigeon/chain/evm/mocks"
	chainmocks "github.com/palomachain/pigeon/chain/mocks"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/internal/queue"
	"github.com/palomachain/pigeon/relayer/mocks"
	"github.com/palomachain/pigeon/testutil"
	timemocks "github.com/palomachain/pigeon/util/time/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEstimateMessages(t *testing.T) {
	ctx := context.Background()
	testcases := []struct {
		name   string
		setup  func(t *testing.T) *Relayer
		expErr error
	}{
		{
			name: "without any processor it does nothing",
			setup: func(t *testing.T) *Relayer {
				pc := mocks.NewPalomaClienter(t)
				pc.On("QueryGetEVMChainInfos", mock.Anything, mock.Anything).Return(nil, nil)
				r := New(
					&config.Config{},
					pc,
					evm.NewFactory(evmmocks.NewPalomaClienter(t)),
					timemocks.NewTime(t),
					Config{},
				)
				return r
			},
		},
		{
			name: "it estimates messages",
			setup: func(t *testing.T) *Relayer {
				keyringPass := "abcd"

				dir := t.TempDir()
				keyring := evm.OpenKeystore(dir)
				acc, err := keyring.NewAccount(keyringPass)
				require.NoError(t, err)
				estimates := []chain.MessageWithEstimate{
					{
						MessageWithSignatures: chain.MessageWithSignatures{
							QueuedMessage: chain.QueuedMessage{ID: 1},
						},
						Estimate:           21_000,
						EstimatedByAddress: "validator",
					},
					{
						MessageWithSignatures: chain.MessageWithSignatures{
							QueuedMessage: chain.QueuedMessage{ID: 2},
						},
						Estimate:           22_000,
						EstimatedByAddress: "validator",
					},
					{
						MessageWithSignatures: chain.MessageWithSignatures{
							QueuedMessage: chain.QueuedMessage{ID: 3},
						},
						Estimate:           23_000,
						EstimatedByAddress: "validator",
					},
				}

				p := chainmocks.NewProcessor(t)
				p.On("IsRightChain", mock.Anything).Return(nil)
				p.On("SupportedQueues").Return([]string{"a"})
				p.On(
					"EstimateMessages",
					mock.Anything,
					queue.FromString("a"),
					[]chain.MessageWithSignatures{
						{QueuedMessage: chain.QueuedMessage{ID: 1}},
						{QueuedMessage: chain.QueuedMessage{ID: 2}},
						{QueuedMessage: chain.QueuedMessage{ID: 3}},
					},
				).Return(estimates, nil)

				pal := mocks.NewPalomaClienter(t)
				pal.On("QueryGetEVMChainInfos", mock.Anything, mock.Anything).Return([]*evmtypes.ChainInfo{
					{
						ChainReferenceID:      "main",
						ChainID:               5,
						SmartContractUniqueID: []byte("5"),
						SmartContractAddr:     common.BytesToAddress([]byte("abcd")).Hex(),
						ReferenceBlockHeight:  5,
						ReferenceBlockHash:    "0x12",
						MinOnChainBalance:     "10000",
					},
				}, nil)
				pal.On("QueryMessagesForEstimating", mock.Anything, mock.Anything).Return(
					[]chain.MessageWithSignatures{
						{QueuedMessage: chain.QueuedMessage{ID: 1}},
						{QueuedMessage: chain.QueuedMessage{ID: 2}},
						{QueuedMessage: chain.QueuedMessage{ID: 3}},
					},
					nil,
				)
				pal.On("AddMessagesGasEstimate", mock.Anything, mock.Anything, estimates[0], estimates[1], estimates[2]).Return(nil)

				os.Setenv("TEST_PASS", keyringPass)
				t.Cleanup(func() {
					os.Unsetenv("TEST_PASS")
				})

				factory := mocks.NewEvmFactorier(t)

				factory.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(p, nil)

				return New(
					&config.Config{
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
					timemocks.NewTime(t),
					Config{},
				)
			},
		},
	}

	for _, tt := range testcases {
		asserter := assert.New(t)
		t.Run(tt.name, func(t *testing.T) {
			relayer := tt.setup(t)

			var locker testutil.FakeMutex
			actualErr := relayer.EstimateMessages(ctx, &locker)
			asserter.Equal(tt.expErr, actualErr)
		})
	}
}
