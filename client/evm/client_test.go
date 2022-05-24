package evm

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vizualni/whoops"
)

func TestReadingStoredContracts(t *testing.T) {
	t.Run("it successfully reads the hello.json contract", func(t *testing.T) {
		c := StoredContracts()
		require.GreaterOrEqual(t, len(c), 1)
		require.Contains(t, c, "hello")
		fmt.Println(c["hello"])
	})
}

func TestExecutingSmartContract(t *testing.T) {
	require.Contains(t, StoredContracts(), "simple")
	ctx := context.Background()
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
			ks.Unlock(acc, "bla")
			args := executeSmartContractIn{
				chainID:       big.NewInt(1337),
				gasAdjustment: 1.0,
				contract:      common.HexToAddress("0xBABA"),
				signingAddr:   acc.Address,
				abi:           StoredContracts()["simple"],
				method:        "store",
				arguments:     []any{big.NewInt(123)},
				keystore:      ks,
			}

			tt.setup(t, &args)

			err = executeSmartContract(
				ctx,
				args,
			)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
