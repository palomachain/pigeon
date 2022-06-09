package relayer

import (
	"context"

	"github.com/palomachain/sparrow/chain"
	"github.com/palomachain/sparrow/chain/paloma"
	"github.com/palomachain/sparrow/util/slice"
)

func (r *Relayer) updateExternalChainInfos(ctx context.Context) error {
	// this returns info about the current validator's keys.
	// if this sparrow is trying to register with a key that the other validator
	// has registered, then it's going to fail as well!
	existingAccInfo, err := r.palomaClient.QueryValidatorInfo(ctx, r.validatorAddress)
	if err != nil {
		return err
	}

	externalAccounts := slice.Map(
		slice.FromMapValues(r.processors),
		func(p chain.Processor) chain.ExternalAccount {
			return p.ExternalAccount()
		},
	)

	chainInfos := []paloma.ChainInfoIn{}

	// TODO: implement accounts removal
	for _, accAddr := range externalAccounts {

		// check if this acc is already registered
		found := false
		for _, currentChainInfo := range existingAccInfo {
			if accAddr.ChainType == currentChainInfo.ChainType &&
				accAddr.ChainID == currentChainInfo.ChainID &&
				accAddr.Address == currentChainInfo.Address {
				found = true
				break
			}
		}
		if !found {
			chainInfos = append(chainInfos, paloma.ChainInfoIn{
				ChainID:    accAddr.ChainID,
				AccAddress: accAddr.Address,
				ChainType:  accAddr.ChainType,
				PubKey:     accAddr.PubKey,
			})
		}
	}

	if len(chainInfos) == 0 {
		return nil
	}

	return r.palomaClient.AddExternalChainInfo(ctx, chainInfos...)
}
