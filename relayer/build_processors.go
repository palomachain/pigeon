package relayer

import (
	"context"

	"github.com/palomachain/sparrow/chain"
	"github.com/palomachain/sparrow/errors"
	"github.com/vizualni/whoops"
)

var (
	ErrMissingChainConfig = errors.Unrecoverable(whoops.String("missing chain config"))
)

func (r *Relayer) buildProcessors(ctx context.Context) ([]chain.Processor, error) {
	chainsInfos, err := r.palomaClient.QueryChainsInfos(ctx)
	if err != nil {
		return err
	}

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

func (r *Relayer) processorFactory(chainInfo chain.ChainInfo) (chain.Processor, error) {
	retErr := whoops.Wrap(ErrMissingChainConfig, whoops.Errorf("reference chain id: %s").Format(chainInfo.ChainReferenceID()))

	switch chainInfo.ChainType() {
	case "EVM":
		cfg, ok := r.config.EVM[chainInfo.ChainReferenceID()]
		if !ok {
			return nil, retErr
		}
		processor, err := r.evmFactory.Build(
			cfg,
			chainInfo.ChainReferenceID(),
		)
		if err != nil {
			return nil, whoops.Wrap(err, retErr)
		}
		return processor, nil
	}

	return nil, retErr
}
