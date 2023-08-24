package evm

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/errors"
	"github.com/palomachain/pigeon/internal/mev"
	arbclient "github.com/roodeag/arbitrum/ethclient"
)

type Factory struct {
	palomaClienter PalomaClienter
}

func NewFactory(pc PalomaClienter) *Factory {
	return &Factory{
		palomaClienter: pc,
	}
}

func (f *Factory) Build(
	cfg config.EVM,
	chainReferenceID,
	smartContractID,
	smartContractABIJson,
	smartContractAddress string,
	chainID *big.Int,
	blockHeight int64,
	blockHeightHash common.Hash,
	minOnChainBalance *big.Int,
	mevClient mev.Client,
) (chain.Processor, error) {
	var smartContractABI *abi.ABI
	if len(smartContractABIJson) > 0 {
		aabi, err := abi.JSON(strings.NewReader(smartContractABIJson))
		if err != nil {
			return Processor{}, errors.Unrecoverable(err)
		}
		smartContractABI = &aabi
	}

	client := &Client{
		config:    cfg,
		paloma:    f.palomaClienter,
		mevClient: mevClient,
	}

	if err := client.init(); err != nil {
		return Processor{}, err
	}

	if chainReferenceID == "arb-main" {
		ac, err := arbclient.Dial(cfg.BaseRPCURL)
		if err != nil {
			return Processor{}, err
		}
		client.arbcon = ac
	}

	return Processor{
		compass: compass{
			CompassID:           smartContractID,
			ChainReferenceID:    chainReferenceID,
			smartContractAddr:   common.HexToAddress(smartContractAddress),
			chainID:             chainID,
			compassAbi:          smartContractABI,
			paloma:              f.palomaClienter,
			evm:                 client,
			startingBlockHeight: blockHeight,
		},
		evmClient:         client,
		chainType:         "evm",
		chainReferenceID:  chainReferenceID,
		blockHeight:       blockHeight,
		blockHeightHash:   blockHeightHash,
		minOnChainBalance: minOnChainBalance,
	}, nil
}
