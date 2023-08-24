package relayer

import (
	"context"
	"sync"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/internal/traits"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) UpdateExternalChainInfos(ctx context.Context, locker sync.Locker) error {
	err := r.buildProcessors(ctx, locker)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("couldn't build processors to update external chain info")

		return err
	}

	log.Info("updating external chain infos")
	externalAccounts := slice.Map(
		r.processors,
		func(p chain.Processor) chain.ExternalAccount {
			return p.ExternalAccount()
		},
	)

	//log.Info("Updating erc20 coin contract for chain")

	chainInfos := slice.Map(externalAccounts, func(acc chain.ExternalAccount) paloma.ChainInfoIn {
		traits := traits.Build(acc.ChainReferenceID, r.mevClient)
		info := paloma.ChainInfoIn{
			ChainReferenceID: acc.ChainReferenceID,
			AccAddress:       acc.Address,
			ChainType:        acc.ChainType,
			PubKey:           acc.PubKey,
			Traits:           traits,
		}
		log.WithFields(log.Fields{
			"chain-reference-id": acc.ChainReferenceID,
			"acc-address":        acc.Address,
			"chain-type":         acc.ChainType,
			"chain-traits":       traits,
		}).Info("sending account info to paloma")
		return info
	})

	if len(chainInfos) == 0 {
		return nil
	}

	locker.Lock()
	err = r.palomaClient.AddExternalChainInfo(ctx, chainInfos...)
	locker.Unlock()

	return err
}
