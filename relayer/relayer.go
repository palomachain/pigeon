package relayer

import (
	"context"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/palomachain/sparrow/attest"
	"github.com/palomachain/sparrow/chain"
	"github.com/palomachain/sparrow/chain/paloma"
	"github.com/palomachain/sparrow/config"
)

type palomaClienter interface {
	KeyName() string
	Keyring() keyring.Keyring
}

type attestExecutor interface {
	Execute(context.Context, string, attest.Request) (attest.Evidence, error)
}

type Relayer struct {
	config config.Root

	// TODO: make an interface for paloma.Client and terra.Client
	palomaClient paloma.Client

	attestExecutor attestExecutor

	validatorAddress sdk.ValAddress

	processors map[string]chain.Processor
}

func New(config config.Root, palomaClient paloma.Client, attestExecutor attestExecutor, processors map[string]chain.Processor) *Relayer {
	return &Relayer{
		config:         config,
		palomaClient:   palomaClient,
		attestExecutor: attestExecutor,
		processors:     processors,
	}
}

func (r *Relayer) init() error {

	r.validatorAddress = r.config.Paloma.ValidatorAddress

	return nil
}
