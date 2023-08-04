package relayer

import (
	"context"
	"math/big"
	"sync"

	"github.com/VolumeFi/whoops"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/palomachain/paloma/x/evm/types"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/errors"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) buildProcessors(ctx context.Context, locker sync.Locker) error {
	locker.Lock()
	defer locker.Unlock()
	queriedChainsInfos, err := r.palomaClient.QueryGetEVMChainInfos(ctx)
	if err != nil {
		return err
	}
	logger := log.WithFields(log.Fields{})

	logger.WithField("chains-infos", queriedChainsInfos).Trace("got chain infos")

	// See if we need to update
	if (r.processors != nil) && (r.chainsInfos != nil) && (len(r.chainsInfos) == len(queriedChainsInfos)) {
		chainsChanged := false
		for k, c := range r.chainsInfos {
			if c.Id != queriedChainsInfos[k].Id ||
				c.ChainReferenceID != queriedChainsInfos[k].ChainReferenceID ||
				c.ChainID != queriedChainsInfos[k].ChainID ||
				string(c.SmartContractUniqueID) != string(queriedChainsInfos[k].SmartContractUniqueID) ||
				c.SmartContractAddr != queriedChainsInfos[k].SmartContractAddr ||
				c.ReferenceBlockHeight != queriedChainsInfos[k].ReferenceBlockHeight ||
				c.ReferenceBlockHash != queriedChainsInfos[k].ReferenceBlockHash ||
				c.Abi != queriedChainsInfos[k].Abi ||
				string(c.Bytecode) != string(queriedChainsInfos[k].Bytecode) ||
				string(c.ConstructorInput) != string(queriedChainsInfos[k].ConstructorInput) ||
				c.Status != queriedChainsInfos[k].Status ||
				c.ActiveSmartContractID != queriedChainsInfos[k].ActiveSmartContractID ||
				c.MinOnChainBalance != queriedChainsInfos[k].MinOnChainBalance {
				chainsChanged = true
			}
		}
		if !chainsChanged {
			logger.Debug("chain infos unchanged since last tick")
			return nil
		}
	}

	logger.Debug("chain infos changed.  building processors")

	r.processors = []chain.Processor{}
	r.chainsInfos = []evmtypes.ChainInfo{}
	for _, chainInfo := range queriedChainsInfos {
		logger = logger.WithFields(log.Fields{
			"chain-reference-id": chainInfo.GetChainReferenceID(),
		})
		processor, err := r.processorFactory(chainInfo)
		if errors.IsUnrecoverable(err) {
			logger.WithError(err).Error("unable to build processor")
			return err
		}

		if err := processor.IsRightChain(ctx); err != nil {
			logger.WithError(err).Error("incorrect chain")
			return err
		}

		r.processors = append(r.processors, processor)
		r.chainsInfos = append(r.chainsInfos, *chainInfo)
	}

	return nil
}

func (r *Relayer) processorFactory(chainInfo *evmtypes.ChainInfo) (chain.Processor, error) {
	// TODO: add support of other types of chains! Right now, only EVM types are supported!
	retErr := whoops.Wrap(ErrMissingChainConfig, whoops.Errorf("reference chain id: %s").Format(chainInfo.GetChainReferenceID()))

	cfg, ok := r.config.EVM[chainInfo.GetChainReferenceID()]
	if !ok {
		return nil, retErr
	}

	chainID := big.NewInt(int64(chainInfo.GetChainID()))

	minOnChainBalance, ok := new(big.Int).SetString(chainInfo.GetMinOnChainBalance(), 10)
	if !ok {
		return nil, ErrInvalidMinOnChainBalance.Format(chainInfo.GetMinOnChainBalance())
	}

	processor, err := r.evmFactory.Build(
		cfg,
		chainInfo.GetChainReferenceID(),
		string(chainInfo.GetSmartContractUniqueID()),
		chainInfo.GetAbi(),
		chainInfo.GetSmartContractAddr(),
		chainID,
		int64(chainInfo.GetReferenceBlockHeight()),
		common.HexToHash(chainInfo.GetReferenceBlockHash()),
		minOnChainBalance,
	)
	if err != nil {
		return nil, whoops.Wrap(err, retErr)
	}
	return processor, nil
}
