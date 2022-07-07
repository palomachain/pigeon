package relayer

import (
	"context"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/util/slice"
)

func (r *Relayer) updateExternalChainInfos(ctx context.Context, processors []chain.Processor) error {
	// this returns info about the current validator's keys.
	// if this pigeon is trying to register with a key that the other validator
	// has registered, then it's going to fail as well!
	existingAccInfo, err := r.palomaClient.QueryValidatorInfo(ctx)
	if err != nil {
		return err
	}

	externalAccounts := slice.Map(
		processors,
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
				accAddr.ChainReferenceID == currentChainInfo.ChainReferenceID &&
				accAddr.Address == currentChainInfo.Address {
				found = true
				break
			}
		}
		if !found {
			chainInfos = append(chainInfos, paloma.ChainInfoIn{
				ChainReferenceID: accAddr.ChainReferenceID,
				AccAddress:       accAddr.Address,
				ChainType:        accAddr.ChainType,
				PubKey:           accAddr.PubKey,
			})
		}
	}

	if len(chainInfos) == 0 {
		return nil
	}

	return r.palomaClient.AddExternalChainInfo(ctx, chainInfos...)
}
