package evm

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/palomachain/sparrow/config"
	"github.com/palomachain/sparrow/errors"
)

type Factory struct {
	palomaClienter palomaClienter
}

func NewFactory(pc palomaClienter) *Factory {
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
) (Processor, error) {

	var smartContractABI *abi.ABI
	if len(smartContractABIJson) > 0 {
		aabi, err := abi.JSON(strings.NewReader(smartContractABIJson))
		if err != nil {
			return Processor{}, errors.Unrecoverable(err)
		}
		smartContractABI = &aabi
	}

	client := &Client{
		config: cfg,
		paloma: f.palomaClienter,
	}

	if err := client.init(); err != nil {
		return Processor{}, err
	}

	// if !ethcommon.IsHexAddress(smartContractAddress) {
	// 	return Processor{}, errors.Unrecoverable(ErrInvalidAddress.Format(smartContractAddress))
	// }

	return Processor{
		compass: compass{
			CompassID:         smartContractID,
			ChainReferenceID:  chainReferenceID,
			smartContractAddr: common.HexToAddress(smartContractAddress),
			chainID:           chainID,
			compassAbi:        smartContractABI,
			paloma:            f.palomaClienter,
			evm:               client,
		},
		evmClient:        client,
		chainType:        "EVM",
		chainReferenceID: chainReferenceID,
	}, nil
}
