package relayer

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/VolumeFi/whoops"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/palomachain/paloma/x/evm/types"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/errors"
	"github.com/palomachain/pigeon/internal/liblog"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) buildProcessors(ctx context.Context, _ sync.Locker) error {
	logger := liblog.WithContext(ctx)
	// TODO: This should live in a short lived cache, it's run very often.
	queriedChainsInfos, err := r.palomaClient.QueryGetEVMChainInfos(ctx)
	if err != nil {
		return err
	}

	// See if we need to update
	r.procRefreshMutex.RLock()
	err = r.validateChainInfos(queriedChainsInfos)
	r.procRefreshMutex.RUnlock()
	if err == nil {
		logger.Debug("chain infos unchanged since last tick")
		return nil
	}

	logger.WithError(err).Warn("Chain infos changed. Building processors...")
	logger.Debug("Acquiring mutex...")
	r.procRefreshMutex.Lock()
	logger.Debug("Mutex acquired.")
	defer func() {
		logger.Debug("Releasing mutex.")
		r.procRefreshMutex.Unlock()
	}()

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

	cfg, ok := r.cfg.EVM[chainInfo.GetChainReferenceID()]
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
		chainInfo.GetFeeManagerAddr(),
		chainID,
		int64(chainInfo.GetReferenceBlockHeight()),
		common.HexToHash(chainInfo.GetReferenceBlockHash()),
		minOnChainBalance,
		r.mevClient,
	)
	if err != nil {
		return nil, whoops.Wrap(err, retErr)
	}
	return processor, nil
}

func (r *Relayer) validateChainInfos(q []*evmtypes.ChainInfo) error {
	if r.processors == nil {
		return fmt.Errorf("missing processors")
	}

	if r.chainsInfos == nil {
		return fmt.Errorf("missing chains infos")
	}

	if err := cmpIfEq("amount of chains", len(r.chainsInfos), len(q)); err != nil {
		return err
	}

	for i, v := range r.chainsInfos {

		type prdct struct {
			name string
			want interface{}
			got  interface{}
		}

		predicate := func(n string, w interface{}, g interface{}) prdct {
			return prdct{name: n, want: w, got: g}
		}

		for _, k := range []prdct{
			predicate("Id", v.Id, q[i].Id),
			predicate("ChainReferenceID", v.ChainReferenceID, q[i].ChainReferenceID),
			predicate("ChainID", v.ChainID, q[i].ChainID),
			predicate("SmartContractUniqueID", string(v.SmartContractUniqueID), string(q[i].SmartContractUniqueID)),
			predicate("SmartContractAddr", v.SmartContractAddr, q[i].SmartContractAddr),
			predicate("ReferenceBlockHeight", v.ReferenceBlockHeight, q[i].ReferenceBlockHeight),
			predicate("ReferenceBlockHash", v.ReferenceBlockHash, q[i].ReferenceBlockHash),
			predicate("Abi", v.Abi, q[i].Abi),
			predicate("Bytecode", string(v.Bytecode), string(q[i].Bytecode)),
			predicate("ConstructorInput", string(v.ConstructorInput), string(q[i].ConstructorInput)),
			predicate("Status", v.Status, q[i].Status),
			predicate("ActiveSmartContractID", v.ActiveSmartContractID, q[i].ActiveSmartContractID),
			predicate("MinOnChainBalance", v.MinOnChainBalance, q[i].MinOnChainBalance),
		} {
			if err := cmpIfEq(k.name, k.want, k.got); err != nil {
				return err
			}
		}
	}

	return nil
}

func cmpIfEq[K comparable](s string, want, got K) error {
	if got == want {
		return nil
	}

	return fmt.Errorf("chain info mismatch in: '%s', want '%v', got '%v'", s, want, got)
}
