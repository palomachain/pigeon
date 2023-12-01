package ethfilter

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type factory struct {
	p         BlockNumberProvider
	addresses []common.Address
	topics    [][]common.Hash
	margin    int64
}
type BlockNumberProvider func(context.Context) (*big.Int, error)

func Factory() *factory { return &factory{} }

func (f *factory) WithFromBlockNumberProvider(provider BlockNumberProvider) *factory {
	f.p = provider
	return f
}

func (f *factory) WithFromBlockNumberSafetyMargin(margin int64) *factory {
	f.margin = margin
	return f
}

func (f *factory) WithAddresses(addresses ...common.Address) *factory {
	f.addresses = addresses
	return f
}

func (f *factory) WithTopics(topics ...[]common.Hash) *factory {
	f.topics = topics
	return f
}

func (f *factory) Filter(context.Context) (ethereum.FilterQuery, error) {
	if f.p == nil {
		return ethereum.FilterQuery{}, errors.New("missing from block number provider")
	}

	b, err := f.p(context.Background())
	if err != nil {
		return ethereum.FilterQuery{}, err
	}

	fromBlock := big.NewInt(0).Sub(b, big.NewInt(f.margin))

	return ethereum.FilterQuery{
		Addresses: f.addresses,
		Topics:    f.topics,
		FromBlock: fromBlock,
	}, nil
}
