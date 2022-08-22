package relayer

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vizualni/whoops"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (r *Relayer) HealthCheck(ctx context.Context) error {
	chainsInfos, err := r.palomaClient.QueryGetEVMChainInfos(ctx)

	if err != nil {
		return err
	}

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
