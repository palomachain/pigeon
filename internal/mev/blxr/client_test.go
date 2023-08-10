package blxr_test

import (
	"testing"

	"github.com/palomachain/pigeon/internal/mev/blxr"
	"github.com/stretchr/testify/require"
)

func TestClientRegisterChain(t *testing.T) {
	c := blxr.New("auth")
	tests := []struct {
		name          string
		chainID       string
		expectedPanic bool
	}{
		{
			name:          "with eth",
			chainID:       "eth-main",
			expectedPanic: false,
		},
		{
			name:          "with bnb",
			chainID:       "bnb-main",
			expectedPanic: false,
		},
		{
			name:          "with matic",
			chainID:       "matic-main",
			expectedPanic: false,
		},
		{
			name:          "with double registry chain",
			chainID:       "eth-main",
			expectedPanic: true,
		},
		{
			name:          "with unsupported chain",
			chainID:       "kava-main",
			expectedPanic: true,
		},
	}

	for _, v := range tests {
		if v.expectedPanic {
			require.Panics(t, func() { c.RegisterChain(v.chainID) }, v.name)
		} else {
			require.NotPanics(t, func() { c.RegisterChain(v.chainID) }, v.name)
		}
	}
}

func TestClientIsChainRegistered(t *testing.T) {
	c := blxr.New("auth")
	tests := []struct {
		name         string
		chainID      string
		expectedFind bool
	}{
		{
			name:         "with eth",
			chainID:      "eth-main",
			expectedFind: true,
		},
		{
			name:         "with bnb",
			chainID:      "bnb-main",
			expectedFind: true,
		},
		{
			name:         "with matic",
			chainID:      "matic-main",
			expectedFind: false,
		},
	}

	c.RegisterChain("eth-main")
	c.RegisterChain("bnb-main")

	for _, v := range tests {
		require.Equal(t, v.expectedFind, c.IsChainRegistered(v.chainID), v.name)
	}

	t.Run("with not supported chain type", func(t *testing.T) {
		require.Panics(t, func() { c.IsChainRegistered("kava-main") })
	})
}
