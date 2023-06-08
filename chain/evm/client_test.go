package evm

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/VolumeFi/whoops"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	log1 = types.Log{
		BlockNumber: 1,
	}
	log2 = types.Log{
		BlockNumber: 2,
	}
	log3 = types.Log{
		BlockNumber: 3,
	}
	log4 = types.Log{
		BlockNumber: 4,
	}
	log5 = types.Log{
		BlockNumber: 5,
	}
)

func TestReadingStoredContracts(t *testing.T) {
	t.Run("it successfully reads the hello.json contract", func(t *testing.T) {
		c := StoredContracts()
		require.GreaterOrEqual(t, len(c), 1)
		require.Contains(t, c, "hello")
	})
}

func TestExecutingSmartContract(t *testing.T) {
	require.Contains(t, StoredContracts(), "simple")
	// ctx := context.Background()
	cryptokey, err := crypto.HexToECDSA(privateKeyBob)
	require.NoError(t, err)
	fakeErr := whoops.String("oh no")
	for _, tt := range []struct {
		name        string
		setup       func(t *testing.T, args *executeSmartContractIn)
		expectedErr error
	}{
		{
			name: "happy path",
			setup: func(t *testing.T, args *executeSmartContractIn) {
				ethMock := newMockEthClienter(t)

				ethMock.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(333), nil)

				ethMock.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(444), nil)

				ethMock.On("SuggestGasTipCap", mock.Anything).Return(big.NewInt(4), nil)

				ethMock.On("PendingCodeAt", mock.Anything, args.contract).Return([]byte("a"), nil)

				ethMock.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(222), nil)

				ethMock.On("SendTransaction", mock.Anything, mock.Anything).Return(nil)

				args.ethClient = ethMock
			},
		},
		{
			name:        "nonce returns an error and it returns error back",
			expectedErr: fakeErr,
			setup: func(t *testing.T, args *executeSmartContractIn) {
				ethMock := newMockEthClienter(t)

				ethMock.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(0), fakeErr)

				args.ethClient = ethMock
			},
		},
		{
			name:        "sending transaction returns and error and that error is sent back",
			expectedErr: fakeErr,
			setup: func(t *testing.T, args *executeSmartContractIn) {
				ethMock := newMockEthClienter(t)

				ethMock.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(333), nil)

				ethMock.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(444), nil)

				ethMock.On("SuggestGasTipCap", mock.Anything).Return(big.NewInt(4), nil)

				ethMock.On("PendingCodeAt", mock.Anything, args.contract).Return([]byte("a"), nil)

				ethMock.On("EstimateGas", mock.Anything, mock.Anything).Return(uint64(222), nil)

				ethMock.On("SendTransaction", mock.Anything, mock.Anything).Return(fakeErr)

				args.ethClient = ethMock
			},
		},
		{
			name:        "gas estimation returns an error and it returns error back",
			expectedErr: fakeErr,
			setup: func(t *testing.T, args *executeSmartContractIn) {
				ethMock := newMockEthClienter(t)

				ethMock.On("PendingNonceAt", mock.Anything, mock.Anything).Return(uint64(333), nil)

				ethMock.On("SuggestGasPrice", mock.Anything).Return(nil, fakeErr)

				args.ethClient = ethMock
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ks := keystore.NewKeyStore(t.TempDir(), keystore.StandardScryptN, keystore.StandardScryptP)
			require.NoError(t, err)
			acc, err := ks.ImportECDSA(cryptokey, "bla")
			require.NoError(t, err)
			ks.Unlock(acc, "bla")
			contract := StoredContracts()["simple"]
			args := executeSmartContractIn{
				chainID:       big.NewInt(1337),
				gasAdjustment: 2.0,
				txType:        2,
				contract:      common.HexToAddress("0xBABA"),
				signingAddr:   acc.Address,
				abi:           contract.ABI,
				method:        "store",
				arguments:     []any{big.NewInt(123)},
				keystore:      ks,
			}

			tt.setup(t, &args)

			_, err = callSmartContract(
				context.Background(),
				args,
			)

			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestFilterLogs(t *testing.T) {
	fakeErr := whoops.String("fake error")

	hashPtr := func(h common.Hash) *common.Hash {
		return &h
	}

	for _, tt := range []struct {
		name string

		filterQuery ethereum.FilterQuery
		blockHeight *big.Int
		callback    func(*testing.T) func([]types.Log) bool

		setup func(t *testing.T) *mockEthClientToFilterLogs

		expErr error
		expRes []types.Log
	}{
		{
			name: "unable to get the current block height returns an error",
			setup: func(t *testing.T) *mockEthClientToFilterLogs {
				srv := newMockEthClientToFilterLogs(t)
				srv.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(nil, fakeErr)
				return srv
			},
			expErr: fakeErr,
		},
		{
			name: "if it's looking at exact block, then it will not change any arguments",
			filterQuery: ethereum.FilterQuery{
				BlockHash: hashPtr(common.HexToHash("abc")),
			},
			setup: func(t *testing.T) *mockEthClientToFilterLogs {
				srv := newMockEthClientToFilterLogs(t)
				srv.On("HeaderByNumber", mock.Anything, mock.Anything).Return(&types.Header{
					Number: big.NewInt(134),
				}, nil)
				srv.On("FilterLogs", mock.Anything, ethereum.FilterQuery{
					BlockHash: hashPtr(common.HexToHash("abc")),
				}).Return(nil, nil)
				return srv
			},
		},
		{
			name:        "if BlockHash is nil, and ToBlock and/or FromBlocks are also nil, then they are changed",
			filterQuery: ethereum.FilterQuery{},
			blockHeight: big.NewInt(134),
			setup: func(t *testing.T) *mockEthClientToFilterLogs {
				srv := newMockEthClientToFilterLogs(t)
				srv.On("FilterLogs", mock.Anything, ethereum.FilterQuery{
					ToBlock:   big.NewInt(134),
					FromBlock: big.NewInt(0),
				}).Return(nil, nil)
				return srv
			},
		},
		{
			name:        "if there are results, then callback is called",
			filterQuery: ethereum.FilterQuery{},
			blockHeight: big.NewInt(134),
			callback: func(t *testing.T) func([]types.Log) bool {
				return func(logs []types.Log) bool {
					require.Empty(t, logs)
					return true
				}
			},
			setup: func(t *testing.T) *mockEthClientToFilterLogs {
				srv := newMockEthClientToFilterLogs(t)
				srv.On("FilterLogs", mock.Anything, ethereum.FilterQuery{
					ToBlock:   big.NewInt(134),
					FromBlock: big.NewInt(0),
				}).Return(nil, nil)
				return srv
			},
		},
		{
			name:        "if there are more than 10000 results, then it calls it recursively",
			filterQuery: ethereum.FilterQuery{},
			blockHeight: big.NewInt(8),
			setup: func(t *testing.T) *mockEthClientToFilterLogs {
				results := [][]types.Log{
					{log4, log5},
					{log1, log2, log3},
				}
				srv := newMockEthClientToFilterLogs(t)

				srv.On("FilterLogs", mock.Anything, ethereum.FilterQuery{
					ToBlock:   big.NewInt(8),
					FromBlock: big.NewInt(0),
				}).Times(1).Return(nil, whoops.String("query returned more than 10000 results"))
				callResults := func(from, to int64, index int) {
					srv.On("FilterLogs", mock.Anything, ethereum.FilterQuery{
						FromBlock: big.NewInt(from),
						ToBlock:   big.NewInt(to),
					}).Times(1).Return(results[index], nil)
				}

				callResults(5, 8, 0)
				callResults(0, 4, 1)
				return srv
			},
			expRes: []types.Log{log5, log4, log3, log2, log1},
		},
		{
			name:        "any other error is returned",
			filterQuery: ethereum.FilterQuery{},
			blockHeight: big.NewInt(8),
			setup: func(t *testing.T) *mockEthClientToFilterLogs {
				srv := newMockEthClientToFilterLogs(t)
				srv.On("FilterLogs", mock.Anything, ethereum.FilterQuery{
					ToBlock:   big.NewInt(8),
					FromBlock: big.NewInt(0),
				}).Return(nil, fakeErr)
				return srv
			},
			expErr: fakeErr,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ethClienter := tt.setup(t)

			var res []types.Log
			defaultCallback := func(logs []types.Log) bool {
				res = append(res, logs...)
				return false
			}

			var fnc func([]types.Log) bool
			if tt.callback != nil {
				fn := tt.callback(t)
				fnc = func(logs []types.Log) bool {
					defaultCallback(logs)
					return fn(logs)
				}
			} else {
				fnc = defaultCallback
			}

			_, err := filterLogs(ctx, ethClienter, tt.filterQuery, tt.blockHeight, true, fnc)

			require.ErrorIs(t, err, tt.expErr)
			require.Equal(t, tt.expRes, res)
		})
	}
}

func TestFindingTheBlockNearestToTime(t *testing.T) {
	type ethHeader struct {
		height uint64
		time   uint64
	}

	testdata := []struct {
		name string

		start uint64
		when  time.Time

		currentBlockNumber uint64
		headers            []ethHeader

		expErr    error
		expHeight uint64
	}{
		{
			name:               "one block fails because the same block is the current one",
			start:              1,
			when:               time.Unix(100, 0),
			currentBlockNumber: 1,
			headers: []ethHeader{
				{
					height: 1,
					time:   1,
				},
			},
			expErr: ErrBlockNotYetGenerated,
		},
		{
			name:               "one block but the time is in the future",
			start:              1,
			when:               time.Unix(100, 0),
			currentBlockNumber: 1,
			headers: []ethHeader{
				{
					height: 1,
					time:   200,
				},
			},
			expErr:    ErrStartingBlockIsInTheFuture,
			expHeight: 0,
		},
		{
			name:               "blocks all in future",
			start:              1,
			when:               time.Unix(100, 0),
			currentBlockNumber: 3,
			headers: []ethHeader{
				{
					height: 1,
					time:   200,
				},
				{
					height: 2,
					time:   300,
				},
				{
					height: 3,
					time:   400,
				},
			},
			expErr:    ErrStartingBlockIsInTheFuture,
			expHeight: 0,
		},
		{
			name:               "block not yet generated",
			start:              1,
			when:               time.Unix(401, 0),
			currentBlockNumber: 3,
			headers: []ethHeader{
				{
					height: 1,
					time:   200,
				},
				{
					height: 2,
					time:   300,
				},
				{
					height: 3,
					time:   400,
				},
			},
			expErr:    ErrBlockNotYetGenerated,
			expHeight: 0,
		},
		{
			name:               "block has been generated (even number of blocks)",
			start:              1,
			when:               time.Unix(401, 0),
			currentBlockNumber: 4,
			headers: []ethHeader{
				{
					height: 1,
					time:   200,
				},
				{
					height: 2,
					time:   300,
				},
				{
					height: 3,
					time:   400,
				},
				{
					height: 4,
					time:   500,
				},
			},
			expHeight: 3,
		},
		{
			name:               "block has been generated (odd number of blocks)",
			start:              1,
			when:               time.Unix(401, 0),
			currentBlockNumber: 5,
			headers: []ethHeader{
				{
					height: 1,
					time:   200,
				},
				{
					height: 2,
					time:   300,
				},
				{
					height: 3,
					time:   400,
				},
				{
					height: 4,
					time:   500,
				},
				{
					height: 5,
					time:   600,
				},
			},
			expHeight: 3,
		},
		{
			name:               "result is the first (odd)",
			start:              1,
			when:               time.Unix(202, 0),
			currentBlockNumber: 5,
			headers: []ethHeader{
				{
					height: 1,
					time:   200,
				},
				{
					height: 2,
					time:   300,
				},
				{
					height: 3,
					time:   400,
				},
				{
					height: 4,
					time:   500,
				},
				{
					height: 5,
					time:   600,
				},
			},
			expHeight: 1,
		},
		{
			name:               "result is the first (even)",
			start:              1,
			when:               time.Unix(202, 0),
			currentBlockNumber: 4,
			headers: []ethHeader{
				{
					height: 1,
					time:   200,
				},
				{
					height: 2,
					time:   300,
				},
				{
					height: 3,
					time:   400,
				},
				{
					height: 4,
					time:   500,
				},
			},
			expHeight: 1,
		},
		{
			name:               "result is the next to last (even)",
			start:              1,
			when:               time.Unix(402, 0),
			currentBlockNumber: 4,
			headers: []ethHeader{
				{
					height: 1,
					time:   200,
				},
				{
					height: 2,
					time:   300,
				},
				{
					height: 3,
					time:   400,
				},
				{
					height: 4,
					time:   500,
				},
			},
			expHeight: 3,
		},
		{
			name:               "result is the next to last (odd)",
			start:              1,
			when:               time.Unix(502, 0),
			currentBlockNumber: 5,
			headers: []ethHeader{
				{
					height: 1,
					time:   200,
				},
				{
					height: 2,
					time:   300,
				},
				{
					height: 3,
					time:   400,
				},
				{
					height: 4,
					time:   500,
				},
				{
					height: 5,
					time:   600,
				},
			},
			expHeight: 4,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{}
			m := newMockEthClientConn(t)
			c.conn = m

			m.On("BlockNumber", mock.Anything).Return(tt.currentBlockNumber, nil).Maybe()
			for _, h := range tt.headers {
				m.On("HeaderByNumber", mock.Anything, big.NewInt(int64(h.height))).Return(&ethtypes.Header{
					Time: h.time,
				}, nil).Maybe()
			}

			ctx := context.Background()
			height, err := c.FindBlockNearestToTime(ctx, tt.start, tt.when)
			require.ErrorIs(t, err, tt.expErr)
			require.Equal(t, tt.expHeight, height)
		})
	}
}
