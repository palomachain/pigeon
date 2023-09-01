package blxr

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/palomachain/pigeon/internal/liblog"
)

var chainIdLkup = map[string]string{
	"1":   "eth-main",
	"56":  "bnb-main",
	"137": "matic-main",
}

type requestParams struct {
	Transaction string
}

type result struct {
	ID    string
	Error struct {
		Code    int
		Message string
	}
	Result struct {
		Txhash string
	}
}

func (c *Client) Relay(ctx context.Context, chainID *big.Int, tx *types.Transaction) (common.Hash, error) {
	nilHash := common.Hash{}

	if !c.IsHealthy() {
		return nilHash, fmt.Errorf("client unhealthy")
	}

	chainReferenceID, fnd := chainIdLkup[chainID.String()]
	if !fnd {
		return nilHash, fmt.Errorf("chain EventNonce %s not supported", chainID.String())
	}

	if !c.IsChainRegistered(chainReferenceID) {
		return nilHash, fmt.Errorf("chain %s not registered", chainReferenceID)
	}

	bz, err := tx.MarshalBinary()
	if err != nil {
		return nilHash, fmt.Errorf("failed to encode raw tx hash: %w", err)
	}

	params := requestParams{
		Transaction: hex.EncodeToString(bz),
	}

	liblog.WithContext(ctx).WithField("raw-tx-hash", params.Transaction).Debug("Relaying raw tx.")

	var r result
	res, err := c.rs.
		R().SetContext(ctx).
		SetHeader("Authorization", c.authHeader).
		SetBody(map[string]interface{}{"method": c.chains[chainReferenceID], "id": "1", "params": params}).
		SetResult(&r).
		Post(cBloXRouteCloudAPIURL)
	if err != nil {
		return nilHash, fmt.Errorf("failed to send request: %w", err)
	}

	if res.StatusCode() != http.StatusOK {
		if len(r.Error.Message) > 0 {
			return nilHash, fmt.Errorf(r.Error.Message)
		}

		return nilHash, fmt.Errorf(string(res.Body()))
	}

	if len(r.Result.Txhash) < 1 {
		return nilHash, fmt.Errorf("failed to retrieve hash: %s", string(res.Body()))
	}

	return common.HexToHash(r.Result.Txhash), nil
}
