package relayer

import (
	"context"
	"math/big"

	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/errors"
	evmtypes "github.com/palomachain/pigeon/types/paloma/x/evm/types"
	log "github.com/sirupsen/logrus"
	"github.com/vizualni/whoops"
)

var (
	ErrMissingChainConfig = errors.Unrecoverable(whoops.String("missing chain config"))
)

func (r *Relayer) buildProcessors(ctx context.Context) ([]chain.Processor, error) {
	chainsInfos, err := r.palomaClient.QueryGetEVMChainInfos(ctx)
	if err != nil {
		return nil, err
	}
	log.WithField("chains-infos", chainsInfos).Info("got chain infos")

	processors := []chain.Processor{}
	for _, chainInfo := range chainsInfos {
		processor, err := r.processorFactory(chainInfo)
		if errors.IsUnrecoverable(err) {
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

	processor, err := r.evmFactory.Build(
		cfg,
		chainInfo.GetChainReferenceID(),
		string(chainInfo.GetSmartContractUniqueID()),
		chainInfo.GetAbi(),
		chainInfo.GetSmartContractAddr(),
		chainID,
	)
	if err != nil {
		return nil, whoops.Wrap(err, retErr)
	}
	return processor, nil
}
