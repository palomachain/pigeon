package evm

import (
	"bytes"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/config"
)

func (p Processor) HealthCheck(ctx context.Context) error {
	var zeroAddr common.Address
	if len(p.evmClient.addr) == 0 || bytes.Equal(p.evmClient.addr.Bytes(), zeroAddr.Bytes()) {
		return chain.ErrMissingAccount.Format(p.chainReferenceID)
	}

	balance, err := p.evmClient.BalanceAt(ctx, p.evmClient.addr, 0)
	if err != nil {
		return err
	}

	cmp := balance.Cmp(p.minOnChainBalance)
	if cmp == -1 || balance.Cmp(big.NewInt(0)) == 0 {
		return chain.ErrAccountBalanceLow.Format(balance, p.evmClient.addr, p.chainReferenceID, p.minOnChainBalance)
	}

	return nil
}

func TestAndVerifyConfig(ctx context.Context, cfg config.EVM) error {
	cli := &Client{
		config: cfg,
	}
	err := cli.init()
	if err != nil {
		return err
	}
	_, err = cli.conn.BlockNumber(ctx)
	if err != nil {
		return err
	}
	return nil
}
