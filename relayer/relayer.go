package relayer

import (
	"github.com/99designs/keyring"
	"github.com/palomachain/sparrow/client/paloma"
	"github.com/palomachain/sparrow/client/terra"
	"github.com/palomachain/sparrow/config"
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
}

func New(config config.Root, palomaClient paloma.Client) *Relayer {
	return &Relayer{
		config:       config,
		palomaClient: palomaClient,
	}
}
