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

var methodLookup = map[string]string{
	"eth-main":   "blxr_private_tx",
	"bnb-main":   "bsc_private_tx",
	"matic-main": "polygon_private_tx",
}

type Client struct {
	authHeader string
	isHealthy  bool
	chains     map[string]string
	rs         *resty.Client
}

func New(authHeader string) *Client {
	return &Client{
		authHeader: authHeader,
		chains:     make(map[string]string),
		rs:         resty.New(),
	}
}

func (c *Client) IsHealthy() bool {
	return c.isHealthy
}

func (c *Client) RegisterChain(chainReferenceID string) {
	method, fnd := methodLookup[chainReferenceID]
	if !fnd {
		panic(fmt.Errorf("chain %s not supported for bloXroute MEV support", chainReferenceID))
	}
	if v, fnd := c.chains[chainReferenceID]; fnd {
		panic(fmt.Errorf("chain %s already has an MEV RPC method %s registered", chainReferenceID, v))
	}
	c.chains[chainReferenceID] = method
}

func (c *Client) IsChainRegistered(chainReferenceID string) bool {
	_, fnd := methodLookup[chainReferenceID]
	if !fnd {
		panic(fmt.Errorf("chain %s not supported for bloXroute MEV support", chainReferenceID))
	}

	_, fnd = c.chains[chainReferenceID]
	return fnd
}
