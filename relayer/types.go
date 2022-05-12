package relayer

import (
	consensus "github.com/palomachain/sparrow/types/paloma/x/consensus/types"
	palomatypes "github.com/palomachain/sparrow/types/paloma/x/consensus/types"
)

type consensusMessageQueueType[T consensus.Signable] string

func (c consensusMessageQueueType[T]) queue() string {
	return string(c)
}

var (
	consensusExecuteSmartContract = consensusMessageQueueType[*palomatypes.SignSmartContractExecute]("consensus_kitica")
	consensusUpdateValset         = consensusMessageQueueType[*palomatypes.SignSmartContractExecute]("update_valset")
)
