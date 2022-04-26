package relayer

import (
	"github.com/99designs/keyring"
	"github.com/volumefi/conductor/client/paloma"
	"github.com/volumefi/conductor/client/terra"
	"github.com/volumefi/conductor/config"
)

type cronchainClienter interface {
	KeyName() string
	Keyring() keyring.Keyring
}

type Relayer struct {
	config config.Root

	// TODO: make an interface for paloma.Client and terra.Client
	palomaClient paloma.Client
	terraClients map[string]terra.Client
}

func New(config config.Root, palomaClient paloma.Client) *Relayer {
	return &Relayer{
		config:       config,
		palomaClient: palomaClient,
	}
}
