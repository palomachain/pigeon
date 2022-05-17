package relayer

import (
	"context"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/palomachain/sparrow/attest"
	"github.com/palomachain/sparrow/client/paloma"
	"github.com/palomachain/sparrow/client/terra"
	"github.com/palomachain/sparrow/config"
	"github.com/palomachain/sparrow/errors"
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
	terraClients map[string]terra.Client

	attestExecutor attestExecutor

	signingKeyAddress string
	validatorAddress  string
}

func New(config config.Root, palomaClient paloma.Client, attestExecutor attestExecutor) *Relayer {
	return &Relayer{
		config:         config,
		palomaClient:   palomaClient,
		attestExecutor: attestExecutor,
	}
}

func (r *Relayer) init() error {

	signingKeyInfo, err := r.palomaClient.Keyring().Key(
		r.config.Paloma.SigningKeyName,
	)
	if err != nil {
		return errors.Unrecoverable(err)
	}

	r.signingKeyAddress, err = r.palomaClient.L.EncodeBech32AccAddr(signingKeyInfo.GetAddress())
	r.validatorAddress = r.config.Paloma.ValidatorAddress

	if err != nil {
		return errors.Unrecoverable(err)
	}

	return nil
}
