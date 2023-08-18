package relayer

import (
	"context"
	"testing"

	"github.com/palomachain/paloma/x/evm/types"
	"github.com/palomachain/pigeon/chain"
	chainmocks "github.com/palomachain/pigeon/chain/mocks"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/relayer/mocks"
	"github.com/palomachain/pigeon/testutil"
	timemocks "github.com/palomachain/pigeon/util/time/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBuildProcessors(t *testing.T) {
	testcases := []struct {
		name        string
		setup       func(t *testing.T) (*Relayer, []chain.Processor, []types.ChainInfo)
		expectedErr error
	}{
		{
			name: "when there are no processors on relayer yet it builds processors",
			setup: func(t *testing.T) (*Relayer, []chain.Processor, []types.ChainInfo) {
				chain1Info := types.ChainInfo{
					Id:                1,
					ChainReferenceID:  "chain-1",
					MinOnChainBalance: "5",
				}
				pc := mocks.NewPalomaClienter(t)
				pc.On(
					"QueryGetEVMChainInfos",
					mock.Anything,
					mock.Anything,
				).Return(
					[]*types.ChainInfo{
						&chain1Info,
					},
					nil,
				)

				processorMock := chainmocks.NewProcessor(t)
				processorMock.On("IsRightChain", mock.Anything).Return(nil)

				evmFactoryMock := mocks.NewEvmFactorier(t)
				evmFactoryMock.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(processorMock, nil)

				r := New(
					config.Root{
						EVM: map[string]config.EVM{
							"chain-1": {},
						},
					},
					pc,
					evmFactoryMock,
					timemocks.NewTime(t),
					Config{},
				)
				r.chainsInfos = []types.ChainInfo{
					chain1Info,
				}

				return r,
					[]chain.Processor{
						processorMock,
					},
					[]types.ChainInfo{
						chain1Info,
					}
			},
		},
		{
			name: "when there are no chainsInfos on relayer yet it builds processors",
			setup: func(t *testing.T) (*Relayer, []chain.Processor, []types.ChainInfo) {
				chain1Info := types.ChainInfo{
					Id:                1,
					ChainReferenceID:  "chain-1",
					MinOnChainBalance: "5",
				}
				pc := mocks.NewPalomaClienter(t)
				pc.On(
					"QueryGetEVMChainInfos",
					mock.Anything,
					mock.Anything,
				).Return(
					[]*types.ChainInfo{
						&chain1Info,
					},
					nil,
				)

				processorMock := chainmocks.NewProcessor(t)
				processorMock.On("IsRightChain", mock.Anything).Return(nil)

				evmFactoryMock := mocks.NewEvmFactorier(t)
				evmFactoryMock.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(processorMock, nil)

				r := New(
					config.Root{
						EVM: map[string]config.EVM{
							"chain-1": {},
						},
					},
					pc,
					evmFactoryMock,
					timemocks.NewTime(t),
					Config{},
				)
				r.processors = []chain.Processor{
					chainmocks.NewProcessor(t),
				}

				return r,
					[]chain.Processor{
						processorMock,
					},
					[]types.ChainInfo{
						chain1Info,
					}
			},
		},
		{
			name: "when the chains lengths are different it builds processors",
			setup: func(t *testing.T) (*Relayer, []chain.Processor, []types.ChainInfo) {
				chain1Info := types.ChainInfo{
					Id:                1,
					ChainReferenceID:  "chain-1",
					MinOnChainBalance: "5",
				}
				chain2Info := types.ChainInfo{
					Id:                2,
					ChainReferenceID:  "chain-2",
					MinOnChainBalance: "5",
				}

				pc := mocks.NewPalomaClienter(t)
				pc.On(
					"QueryGetEVMChainInfos",
					mock.Anything,
					mock.Anything,
				).Return(
					[]*types.ChainInfo{
						&chain1Info,
					},
					nil,
				)

				processorMock := chainmocks.NewProcessor(t)
				processorMock.On("IsRightChain", mock.Anything).Return(nil)

				evmFactoryMock := mocks.NewEvmFactorier(t)
				evmFactoryMock.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(processorMock, nil)

				r := New(
					config.Root{
						EVM: map[string]config.EVM{
							"chain-1": {},
						},
					},
					pc,
					evmFactoryMock,
					timemocks.NewTime(t),
					Config{},
				)
				r.processors = []chain.Processor{
					chainmocks.NewProcessor(t),
				}

				r.chainsInfos = []types.ChainInfo{
					chain1Info,
					chain2Info,
				}

				return r,
					[]chain.Processor{
						processorMock,
					},
					[]types.ChainInfo{
						chain1Info,
					}
			},
		},
		{
			name: "when there is a difference in the chain data it builds processors",
			setup: func(t *testing.T) (*Relayer, []chain.Processor, []types.ChainInfo) {
				chain1Info := types.ChainInfo{
					Id:                1,
					ChainReferenceID:  "chain-1",
					MinOnChainBalance: "5",
				}

				chain1NewInfo := types.ChainInfo{
					Id:                1,
					ChainReferenceID:  "chain-1",
					MinOnChainBalance: "50",
				}
				pc := mocks.NewPalomaClienter(t)
				pc.On(
					"QueryGetEVMChainInfos",
					mock.Anything,
					mock.Anything,
				).Return(
					[]*types.ChainInfo{
						&chain1NewInfo,
					},
					nil,
				)

				processorMock := chainmocks.NewProcessor(t)
				processorMock.On("IsRightChain", mock.Anything).Return(nil)

				evmFactoryMock := mocks.NewEvmFactorier(t)
				evmFactoryMock.On("Build", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(processorMock, nil)

				r := New(
					config.Root{
						EVM: map[string]config.EVM{
							"chain-1": {},
						},
					},
					pc,
					evmFactoryMock,
					timemocks.NewTime(t),
					Config{},
				)
				r.processors = []chain.Processor{
					chainmocks.NewProcessor(t),
				}

				r.chainsInfos = []types.ChainInfo{
					chain1Info,
				}

				return r,
					[]chain.Processor{
						processorMock,
					},
					[]types.ChainInfo{
						chain1NewInfo,
					}
			},
		},
		{
			name: "when the chains are the same it doesn't build processors",
			setup: func(t *testing.T) (*Relayer, []chain.Processor, []types.ChainInfo) {
				chain1Info := types.ChainInfo{
					Id:                1,
					ChainReferenceID:  "chain-1",
					MinOnChainBalance: "5",
				}

				pc := mocks.NewPalomaClienter(t)
				pc.On(
					"QueryGetEVMChainInfos",
					mock.Anything,
					mock.Anything,
				).Return(
					[]*types.ChainInfo{
						&chain1Info,
					},
					nil,
				)

				r := New(
					config.Root{
						EVM: map[string]config.EVM{
							"chain-1": {},
						},
					},
					pc,
					mocks.NewEvmFactorier(t),
					timemocks.NewTime(t),
					Config{},
				)

				origProcessor := chainmocks.NewProcessor(t)
				r.processors = []chain.Processor{
					origProcessor,
				}

				r.chainsInfos = []types.ChainInfo{
					chain1Info,
				}

				return r,
					[]chain.Processor{
						origProcessor,
					},
					[]types.ChainInfo{
						chain1Info,
					}
			},
		},
	}

	asserter := assert.New(t)
	ctx := context.Background()
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			relayer, expectedProcessors, expectedChainsInfos := tt.setup(t)
			var locker testutil.FakeMutex

			actualErr := relayer.buildProcessors(ctx, locker)
			asserter.Equal(tt.expectedErr, actualErr)
			asserter.Equal(expectedProcessors, relayer.processors)
			asserter.Equal(expectedChainsInfos, relayer.chainsInfos)
		})
	}
}
