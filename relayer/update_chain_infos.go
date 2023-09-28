package relayer

import (
	"context"
	"reflect"
	"sync"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/internal/liblog"
	"github.com/palomachain/pigeon/internal/traits"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) UpdateExternalChainInfos(ctx context.Context, locker sync.Locker) error {
	logger := liblog.WithContext(ctx).WithField("component", "update-external-chain-infos")
	logger.Info("updating external chain infos")

	if err := r.buildProcessors(ctx, locker); err != nil {
		logger.WithError(err).Error("couldn't build processors to update external chain info")
		return err
	}

	externalAccounts := make([]chain.ExternalAccount, len(r.processors))
	for i, v := range r.processors {
		externalAccounts[i] = v.ExternalAccount()
	}

	chainInfos := make([]paloma.ChainInfoIn, len(externalAccounts))
	for i, v := range externalAccounts {
		traits := traits.Build(v.ChainReferenceID, r.mevClient)
		logger.WithFields(log.Fields{
			"chain-reference-id": v.ChainReferenceID,
			"acc-address":        v.Address,
			"chain-type":         v.ChainType,
			"chain-traits":       traits,
		}).Info("adding chain info to payload")

		chainInfos[i] = paloma.ChainInfoIn{
			ChainReferenceID: v.ChainReferenceID,
			AccAddress:       v.Address,
			ChainType:        v.ChainType,
			PubKey:           v.PubKey,
			Traits:           traits,
		}
	}

	if len(chainInfos) < 1 {
		return nil
	}

	if reflect.DeepEqual(chainInfos, r.valCache.lastChainInfoRecord) {
		logger.Info("Chain infos unchanged, skip sending...")
		return nil
	}

	locker.Lock()
	err := r.palomaClient.AddExternalChainInfo(ctx, chainInfos...)
	locker.Unlock()

	if err != nil {
		return err
	}

	logger.Info("Updated chain infos record sent, refreshing cache...")
	r.valCache.lastChainInfoRecord = chainInfos
	return nil
}
