package relayer

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/palomachain/paloma/x/evm/types"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/evm"
	evmmocks "github.com/palomachain/pigeon/chain/evm/mocks"
	chainmocks "github.com/palomachain/pigeon/chain/mocks"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/relayer/mocks"
	"github.com/palomachain/pigeon/testutil"
	timemocks "github.com/palomachain/pigeon/util/time/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateExternalChainInfos(t *testing.T) {
	ctx := context.Background()
	locker := testutil.FakeMutex{}
	testErr := fmt.Errorf("Your error is in a different castle!")

	testcases := []struct {
		setup  func(t *testing.T) *Relayer
		expErr error
		name   string
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
			expErr: nil,
		},
		{
			name: "without error during processor build, it returns an error",
			setup: func(t *testing.T) *Relayer {
				pc := mocks.NewPalomaClienter(t)
				pc.On("QueryGetEVMChainInfos", mock.Anything, mock.Anything).Return(nil, testErr)
				r := New(
					&config.Config{},
					pc,
					evm.NewFactory(evmmocks.NewPalomaClienter(t)),
					timemocks.NewTime(t),
					Config{},
				)
				return r
			},
			expErr: testErr,
		},
		{
			name: "without cached values, it will report to Paloma",
			setup: func(t *testing.T) *Relayer {
				keyringPass := "abcd"

				dir := t.TempDir()
				keyring := evm.OpenKeystore(dir)
				acc, err := keyring.NewAccount(keyringPass)
				require.NoError(t, err)

				p := chainmocks.NewProcessor(t)
				p.On("IsRightChain", mock.Anything).Return(nil)
				p.On("ExternalAccount", mock.Anything).Return(chain.ExternalAccount{
					ChainType:        "evm",
					ChainReferenceID: "main",
					Address:          common.BytesToAddress([]byte("abcd")).Hex(),
					PubKey:           acc.Address.Bytes(),
				})

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
				pal.On("AddExternalChainInfo", mock.Anything, mock.Anything).Return(nil)

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
			expErr: nil,
		},
		{
			name: "with cached values, it will skip the call to Paloma",
			setup: func(t *testing.T) *Relayer {
				keyringPass := "abcd"

				dir := t.TempDir()
				keyring := evm.OpenKeystore(dir)
				acc, err := keyring.NewAccount(keyringPass)
				require.NoError(t, err)

				p := chainmocks.NewProcessor(t)
				p.On("IsRightChain", mock.Anything).Return(nil)
				chainAcc := chain.ExternalAccount{
					ChainType:        "evm",
					ChainReferenceID: "main",
					Address:          common.BytesToAddress([]byte("abcd")).Hex(),
					PubKey:           acc.Address.Bytes(),
				}
				p.On("ExternalAccount", mock.Anything).Return(chainAcc)

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

				os.Setenv("TEST_PASS", keyringPass)
				t.Cleanup(func() {
					os.Unsetenv("TEST_PASS")
				})

				factory := mocks.NewEvmFactorier(t)
				factory.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(p, nil)

				r := New(
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

				r.valCache.lastChainInfoRecord = []paloma.ChainInfoIn{{
					ChainType:        "evm",
					ChainReferenceID: chainAcc.ChainReferenceID,
					AccAddress:       chainAcc.Address,
					PubKey:           chainAcc.PubKey,
					Traits:           []string{},
				}}

				return r
			},
			expErr: nil,
		},
		{
			name: "with cached values mismatch, it will make the call to Paloma",
			setup: func(t *testing.T) *Relayer {
				keyringPass := "abcd"

				dir := t.TempDir()
				keyring := evm.OpenKeystore(dir)
				acc, err := keyring.NewAccount(keyringPass)
				require.NoError(t, err)

				p := chainmocks.NewProcessor(t)
				p.On("IsRightChain", mock.Anything).Return(nil)
				chainAcc := chain.ExternalAccount{
					ChainType:        "evm",
					ChainReferenceID: "main",
					Address:          common.BytesToAddress([]byte("abcd")).Hex(),
					PubKey:           acc.Address.Bytes(),
				}
				p.On("ExternalAccount", mock.Anything).Return(chainAcc)

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
				pal.On("AddExternalChainInfo", mock.Anything, mock.Anything).Return(nil)

				os.Setenv("TEST_PASS", keyringPass)
				t.Cleanup(func() {
					os.Unsetenv("TEST_PASS")
				})

				factory := mocks.NewEvmFactorier(t)
				factory.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(p, nil)

				r := New(
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

				r.valCache.lastChainInfoRecord = []paloma.ChainInfoIn{{
					ChainType:        "evm",
					ChainReferenceID: chainAcc.ChainReferenceID,
					AccAddress:       chainAcc.Address,
					PubKey:           chainAcc.PubKey,
					Traits:           []string{"mev"},
				}}

				return r
			},
			expErr: nil,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			relayer := tt.setup(t)
			err := relayer.UpdateExternalChainInfos(ctx, locker)
			require.Equal(t, err, err)
		})
	}
}
