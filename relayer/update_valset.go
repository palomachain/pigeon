package relayer

import (
	"context"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) updateExternalChainInfos(ctx context.Context, processors []chain.Processor) error {
	log.Info("updating external chain infos")
	externalAccounts := slice.Map(
		processors,
		func(p chain.Processor) chain.ExternalAccount {
			return p.ExternalAccount()
		},
	)
	chainInfos := slice.Map(externalAccounts, func(acc chain.ExternalAccount) paloma.ChainInfoIn {
		info := paloma.ChainInfoIn{
			ChainReferenceID: acc.ChainReferenceID,
			AccAddress:       acc.Address,
			ChainType:        acc.ChainType,
			PubKey:           acc.PubKey,
		}
		log.WithFields(log.Fields{
			"ChainReferenceID": acc.ChainReferenceID,
			"AccAddress":       acc.Address,
			"ChainType":        acc.ChainType,
		}).Info("sending account info to paloma")
		return info
	})

	if len(chainInfos) == 0 {
		return nil
	}

	return r.palomaClient.AddExternalChainInfo(ctx, chainInfos...)
}
