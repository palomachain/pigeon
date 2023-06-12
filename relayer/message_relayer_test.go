package relayer

import (
	"context"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/palomachain/paloma/x/evm/types"
	valsettypes "github.com/palomachain/paloma/x/valset/types"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/evm"
	chainmocks "github.com/palomachain/pigeon/chain/mocks"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/relayer/mocks"
	"github.com/palomachain/pigeon/testutil"
	timemocks "github.com/palomachain/pigeon/util/time/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRelayMessages(t *testing.T) {
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
					config.Root{},
					pc,
					evm.NewFactory(pc),
					timemocks.NewTime(t),
					Config{},
				)
				return r
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
				p.On("IsRightChain", mock.Anything).Return(nil)
				p.On("SupportedQueues").Return([]string{"a"})
				p.On(
					"ProcessMessages",
					mock.Anything,
					"a",
					[]chain.MessageWithSignatures{
						{QueuedMessage: chain.QueuedMessage{ID: 123}},
						{QueuedMessage: chain.QueuedMessage{ID: 456}},
					},
				).Return(nil)

				pal := mocks.NewPalomaClienter(t)
				pal.On("GetValidatorAddress").Return(sdk.ValAddress("abc"))
				pal.On("BlockHeight", mock.Anything).Return(int64(555), nil)
				pal.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(
					&valsettypes.Snapshot{
						Validators: []valsettypes.Validator{
							{Address: sdk.ValAddress("abc")},
						},
					},
					nil,
				)
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
				pal.On("QueryMessagesInQueue", mock.Anything, mock.Anything).Return(
					[]chain.MessageWithSignatures{
						{QueuedMessage: chain.QueuedMessage{ID: 123}},
						{QueuedMessage: chain.QueuedMessage{ID: 456}},
						{QueuedMessage: chain.QueuedMessage{ID: 789, PublicAccessData: []byte("tx hash")}},
					},
					nil,
				)

				os.Setenv("TEST_PASS", keyringPass)
				t.Cleanup(func() {
					os.Unsetenv("TEST_PASS")
				})

				factory := mocks.NewEvmFactorier(t)

				factory.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(p, nil)

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
					timemocks.NewTime(t),
					Config{},
				)
			},
		},
		{
			name: "if the processor is connected to the wrong chain it returns the error",
			setup: func(t *testing.T) *Relayer {
				keyringPass := "abcd"

				dir := t.TempDir()
				keyring := evm.OpenKeystore(dir)
				acc, err := keyring.NewAccount(keyringPass)
				require.NoError(t, err)

				p := chainmocks.NewProcessor(t)
				p.On("IsRightChain", mock.Anything).Return(chain.ErrNotConnectedToRightChain)

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

				factory := mocks.NewEvmFactorier(t)

				factory.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(p, nil)

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
					timemocks.NewTime(t),
					Config{},
				)
			},
			expErr: chain.ErrNotConnectedToRightChain,
		},
	}

	for _, tt := range testcases {
		asserter := assert.New(t)
		t.Run(tt.name, func(t *testing.T) {
			relayer := tt.setup(t)

			var locker testutil.FakeMutex
			actualErr := relayer.RelayMessages(ctx, &locker)
			asserter.Equal(tt.expErr, actualErr)
		})
	}
}
