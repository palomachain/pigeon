package relayer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/errors"
	evmtypes "github.com/palomachain/pigeon/types/paloma/x/evm/types"
	log "github.com/sirupsen/logrus"
	"github.com/vizualni/whoops"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (r *Relayer) buildProcessors(ctx context.Context) ([]chain.Processor, error) {
	chainsInfos, err := r.palomaClient.QueryGetEVMChainInfos(ctx)
	if err != nil {
		return nil, err
	}
	log.WithField("chains-infos", chainsInfos).Trace("got chain infos")

	processors := []chain.Processor{}
	for _, chainInfo := range chainsInfos {
		processor, err := r.processorFactory(chainInfo)
		if errors.IsUnrecoverable(err) {
			return nil, err
		}

		if err := processor.IsRightChain(ctx); err != nil {
			return nil, err
		}

		processors = append(processors, processor)
	}

	return processors, nil
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

func (r *Relayer) HealthCheck(ctx context.Context) error {
	chainsInfos, err := r.palomaClient.QueryGetEVMChainInfos(ctx)

	if err != nil {
		return err
	}

	val, err := r.palomaClient.GetValidator(ctx)
	if err != nil {
		return err
	}

	if val == nil {
		// I am sorry, but where *is* your validator?
		return ErrNotAValidatorAccount
	}

	isStaking := false
	if !val.Jailed {
		if val.Status == stakingtypes.Bonded || val.Status == stakingtypes.Unbonding {
			isStaking = true
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
