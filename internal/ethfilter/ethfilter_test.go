package ethfilter_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/palomachain/pigeon/internal/ethfilter"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	ctx := context.Background()

	mockProvider := func(i int64) func(context.Context) (*big.Int, error) {
		return func(context.Context) (*big.Int, error) {
			return big.NewInt(i), nil
		}
	}

	t.Run("with missing block provider", func(t *testing.T) {
		_, err := ethfilter.Factory().Filter(ctx)
		require.Error(t, err)
	})

	t.Run("with failing block provider", func(t *testing.T) {
		p := func(context.Context) (*big.Int, error) { return nil, fmt.Errorf("fail") }
		_, err := ethfilter.Factory().WithFromBlockNumberProvider(p).Filter(ctx)
		require.Error(t, err)
	})

	t.Run("with valid block provider", func(t *testing.T) {
		for _, tt := range []struct {
			name      string
			addresses []common.Address
			topics    [][]common.Hash
			margin    int64
		}{
			{
				name: "with no additional configuration",
			},
			{
				name:      "with addresses",
				addresses: []common.Address{common.HexToAddress("address-1")},
			},
			{
				name:   "with topics",
				topics: [][]common.Hash{{common.HexToHash("topic-1")}},
			},
			{
				name:   "with margin",
				margin: 10,
			},
			{
				name:      "with addresses and topics",
				addresses: []common.Address{common.HexToAddress("address-1")},
				topics:    [][]common.Hash{{common.HexToHash("topic-1")}},
			},
			{
				name:      "with all configuration present",
				addresses: []common.Address{common.HexToAddress("address-1")},
				topics:    [][]common.Hash{{common.HexToHash("topic-1")}},
				margin:    10,
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				p := mockProvider(100)
				f, err := ethfilter.Factory().
					WithFromBlockNumberProvider(p).
					WithAddresses(tt.addresses...).
					WithTopics(tt.topics...).
					WithFromBlockNumberSafetyMargin(tt.margin).
					Filter(ctx)
				require.NoError(t, err)
				require.Equal(t, tt.addresses, f.Addresses)
				require.Equal(t, tt.topics, f.Topics)
				require.Equal(t, int64(100-tt.margin), f.FromBlock.Int64())
			})
		}
	})
}
