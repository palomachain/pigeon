package traits

import (
	valsettypes "github.com/palomachain/paloma/x/valset/types"
)

type Traits []string

type MevRelayerQuery interface {
	IsHealthy() bool
	IsChainRegistered(string) bool
}

func Build(chainID string, mevClient MevRelayerQuery) Traits {
	traits := []string{}

	if mevClient == nil || !mevClient.IsHealthy() || !mevClient.IsChainRegistered(chainID) {
		return traits
	}

	traits = append(traits, valsettypes.PIGEON_TRAIT_MEV)
	return traits
}
