package relayer

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/palomachain/sparrow/client/paloma"
	"github.com/palomachain/sparrow/client/terra"
	"github.com/palomachain/sparrow/config"
	"github.com/palomachain/sparrow/errors"
)

type palomaClienter interface {
	KeyName() string
	Keyring() keyring.Keyring
}

type Relayer struct {
	config config.Root

	// TODO: make an interface for paloma.Client and terra.Client
	palomaClient paloma.Client
	terraClients map[string]terra.Client

	signingKeyAddress string
	validatorAddress  string
}

func New(config config.Root, palomaClient paloma.Client) *Relayer {
	return &Relayer{
		config:       config,
		palomaClient: palomaClient,
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
