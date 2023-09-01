package blxr

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestClientRelay(t *testing.T) {
	authHeader := "foo"
	ctx := context.Background()
	chainID := big.NewInt(1) // eth-main
	chainReferenceID := "eth-main"

	newClient := func(h bool, chains ...string) *Client {
		httpmock.Reset()
		c := New(authHeader)
		c.isHealthy = h
		for _, v := range chains {
			c.RegisterChain(v)
		}
		return c
	}

	newTx := func(legacy bool) *types.Transaction {
		nonce := uint64(20)
		gas := uint64(10000)
		to := common.HexToAddress("e3cd54d29cbf35648edcf53d6a344bd4b88da059")
		value := big.NewInt(1)
		if legacy {
			return types.NewTx(&types.LegacyTx{
				Nonce:    nonce,
				GasPrice: big.NewInt(4),
				Gas:      gas,
				To:       &to,
				Value:    value,
				Data:     []byte{},
			})
		}

		return types.NewTx(&types.DynamicFeeTx{
			ChainID:   big.NewInt(1),
			Nonce:     nonce,
			GasTipCap: big.NewInt(10000),
			GasFeeCap: big.NewInt(10000),
			Gas:       gas,
			To:        &to,
			Value:     value,
			Data:      []byte{},
		})
	}

	t.Run("with unhealthy client", func(t *testing.T) {
		c := newClient(false)
		_, err := c.Relay(ctx, chainID, nil)
		require.Error(t, err, "should return error")
		require.EqualError(t, err, "client unhealthy")
	})

	t.Run("with invalid chain EventNonce", func(t *testing.T) {
		c := newClient(true)
		_, err := c.Relay(ctx, big.NewInt(42), nil)
		require.Error(t, err, "should return error")
		require.EqualError(t, err, "chain EventNonce 42 not supported")
	})

	t.Run("with unregistered chain EventNonce", func(t *testing.T) {
		c := newClient(true)
		_, err := c.Relay(ctx, chainID, nil)
		require.Error(t, err, "should return error")
		require.EqualError(t, err, "chain eth-main not registered")
	})

	t.Run("with unregistered chain EventNonce", func(t *testing.T) {
		c := newClient(true)
		_, err := c.Relay(ctx, chainID, nil)
		require.Error(t, err, "should return error")
		require.EqualError(t, err, "chain eth-main not registered")
	})

	t.Run("with status code not 200", func(t *testing.T) {
		c := newClient(true, chainReferenceID)
		fixture := `{"id":"1","error":{"message":"fail"},"jsonrpc":"2.0"}`
		responder := httpmock.NewStringResponder(400, fixture)
		httpmock.RegisterResponder("POST", cBloXRouteCloudAPIURL, responder)
		httpmock.ActivateNonDefault(c.rs.GetClient())

		_, err := c.Relay(ctx, chainID, types.NewTx(&types.DynamicFeeTx{}))
		require.Error(t, err, "should return error")
		count := httpmock.GetTotalCallCount()
		require.Equal(t, 1, count)
	})

	t.Run("with missing response tx hash", func(t *testing.T) {
		c := newClient(true, chainReferenceID)
		fixture := `{"id":"1","result":{},"jsonrpc":"2.0"}`
		responder := httpmock.NewStringResponder(200, fixture)
		httpmock.RegisterResponder("POST", cBloXRouteCloudAPIURL, responder)
		httpmock.ActivateNonDefault(c.rs.GetClient())

		_, err := c.Relay(ctx, chainID, types.NewTx(&types.DynamicFeeTx{}))
		require.Error(t, err, "should return error")
		require.EqualError(t, err, fmt.Sprintf("failed to retrieve hash: %s", fixture))
		count := httpmock.GetTotalCallCount()
		require.Equal(t, 1, count)
	})

	for name, tx := range map[string]*types.Transaction{
		"with legacy transaction":      newTx(true),
		"with dynamic fee transaction": newTx(false),
	} {
		t.Run(name, func(t *testing.T) {
			c := newClient(true, chainReferenceID)
			fixture := fmt.Sprintf(`{"id":"1","result":{"txHash":"%s"},"jsonrpc":"2.0"}`, hex.EncodeToString(tx.Hash().Bytes()))
			rawTx, err := tx.MarshalBinary()
			require.NoError(t, err)
			expectedBody := fmt.Sprintf(`{"id":"1","method":"blxr_private_tx","params":{"Transaction":"%x"}}`, rawTx)
			responder := httpmock.NewStringResponder(200, fixture).HeaderAdd(map[string][]string{"Content-Type": {"application/json"}}).Once()
			httpmock.RegisterMatcherResponder("POST", cBloXRouteCloudAPIURL, httpmock.BodyContainsString(expectedBody), responder)
			httpmock.ActivateNonDefault(c.rs.GetClient())

			res, err := c.Relay(ctx, chainID, tx)
			require.NoError(t, err)
			require.Equal(t, tx.Hash(), res)
			count := httpmock.GetTotalCallCount()
			require.Equal(t, 1, count)
		})
	}

	httpmock.DeactivateAndReset()
}
