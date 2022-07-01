package evm

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
) (Processor, error) {
	smartContractABI, err := abi.JSON(strings.NewReader(smartContractABIJson))
	if err != nil {
		return Processor{}, errors.Unrecoverable(err)
	}
	client := NewClient(cfg, f.palomaClienter)
	return NewProcessor(client, chainReferenceID, smartContractID, smartContractABI, smartContractAddress), nil
}
