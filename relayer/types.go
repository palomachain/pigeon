package relayer

import (
	"context"
	"github.com/palomachain/sparrow/chain/paloma"
	consensus "github.com/palomachain/sparrow/types/paloma/x/consensus/types"
	palomatypes "github.com/palomachain/sparrow/types/paloma/x/consensus/types"
)

type consensusMessageQueueType[T consensus.Signable] string

func (c consensusMessageQueueType[T]) queue() string {
	return string(c)
}

var (
	consensusExecuteSmartContract = consensusMessageQueueType[*palomatypes.SignSmartContractExecute]("execute_smart_contract")
	consensusUpdateValset         = consensusMessageQueueType[*palomatypes.SignSmartContractExecute]("update_valset")
)

func (c consensusMessageQueueType[T]) relayConsensusReachedMessages(
	ctx context.Context,
	r *Relayer,
) ([]paloma.BroadcastMessageSignatureIn, error) {
	return nil, nil
}
