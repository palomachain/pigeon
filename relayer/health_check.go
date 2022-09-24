package relayer

import (
	"context"
	"errors"
	"strings"

	"github.com/palomachain/pigeon/chain/evm"
	log "github.com/sirupsen/logrus"
	"github.com/vizualni/whoops"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (r *Relayer) isStaking(ctx context.Context) error {

	val, err := r.palomaClient.GetValidator(ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "NotFound") {
			return err
		}
	}

	isStaking := false
	if val != nil {
		if !val.Jailed {
			if val.Status == stakingtypes.Bonded || val.Status == stakingtypes.Unbonding {
				isStaking = true
			}
		}
	}

	if !isStaking {
		return ErrValidatorIsNotStaking
	}
	return nil
}

func (r *Relayer) HealthCheck(ctx context.Context) error {
	chainsInfos, err := r.palomaClient.QueryGetEVMChainInfos(ctx)

	if err != nil {
		return err
	}

	err = r.isStaking(ctx)
	isStaking := false
	if err == nil {
		isStaking = true
	} else {
		if !errors.Is(err, ErrValidatorIsNotStaking) {
			return err
		}
	}

	var g whoops.Group
	for _, chainInfo := range chainsInfos {
		p, err := r.processorFactory(chainInfo)
		if err != nil {
			g.Add(err)
			continue
		}

		g.Add(p.HealthCheck(ctx))
	}

	if !isStaking {
		// then these errors are only warning
		log.Warn("validator is not staking. ensure to fix these warning if you wish to stake.")
		for _, err := range g {
			log.WithError(err).Warn("blocker for becoming a staking validator. Fix if you wish to stake.")
		}

		return nil
	}

	return g.Return()
}

func (r *Relayer) BootHealthCheck(ctx context.Context) error {
	var g whoops.Group
	for _, cfg := range r.config.EVM {
		g.Add(evm.TestAndVerifyConfig(ctx, cfg))
	}
	return g.Return()
}
