package blxr

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	cHealthprobeQueryInterval = time.Second
	cBloXRouteCloudAPIURL     = "https://api.blxrbdn.com"
)

var chainIDLookup = map[string]string{
	"eth-main":   "Mainnet",
	"bnb-main":   "BSC-Mainnet",
	"matic-main": "Polygon-Mainnet",
}

type Client struct {
	authHeader string
	isHealthy  bool
	chains     map[string]struct{}
	rs         *resty.Client
}

func New(authHeader string) *Client {
	return &Client{
		authHeader: authHeader,
		chains:     make(map[string]struct{}),
		rs:         resty.New(),
	}
}

func (c *Client) IsHealthy() bool {
	return c.isHealthy
}

func (c *Client) RegisterChain(referenceChainID string) {
	blxrID, fnd := chainIDLookup[referenceChainID]
	if !fnd {
		panic(fmt.Errorf("chain %s not supported for bloXroute MEV support", referenceChainID))
	}
	if _, fnd := c.chains[blxrID]; fnd {
		panic(fmt.Errorf("chain %s already has an MEV RPC endpoint registered", blxrID))
	}
	c.chains[blxrID] = struct{}{}
}

func (c *Client) IsChainRegistered(referenceChainID string) bool {
	blxrID, fnd := chainIDLookup[referenceChainID]
	if !fnd {
		panic(fmt.Errorf("chain %s not supported for bloXroute MEV support", referenceChainID))
	}

	_, fnd = c.chains[blxrID]
	return fnd
}
