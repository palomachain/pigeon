package relayer

import (
	"github.com/99designs/keyring"
	"github.com/volumefi/conductor/client/cronchain"
	"github.com/volumefi/conductor/client/terra"
)

type cronchainClienter interface {
	KeyName() string
	Keyring() keyring.Keyring
}

type relayer struct {
	palomaClient cronchain.Client
	terraClients map[string]terra.Client
}
