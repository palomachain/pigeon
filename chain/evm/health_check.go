package evm

import (
	"context"

	"github.com/palomachain/pigeon/chain"
)

func (p Processor) HealthCheck(ctx context.Context) error {
	if len(p.evmClient.addr) == 0 {
		return chain.ErrMissingAccount.Format(p.chainReferenceID)
	}
	balance, err := p.evmClient.BalanceAt(ctx, p.evmClient.addr, 0)
	if err != nil {
		return err
	}

	if balance.Cmp(p.minOnChainBalance) == -1 {
		return chain.ErrAccountBalanceLow.Format(balance, p.evmClient.addr, p.chainReferenceID, p.minOnChainBalance)
	}

	return nil
}
