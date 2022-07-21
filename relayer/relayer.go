package relayer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/config"
	evmtypes "github.com/palomachain/pigeon/types/paloma/x/evm/types"
	valset "github.com/palomachain/pigeon/types/paloma/x/valset/types"
)

//go:generate mockery --name=PalomaClienter
type PalomaClienter interface {
	AddExternalChainInfo(ctx context.Context, chainInfos ...paloma.ChainInfoIn) error
	QueryValidatorInfo(ctx context.Context) ([]*valset.ExternalChainInfo, error)
	BroadcastMessageSignatures(ctx context.Context, signatures ...paloma.BroadcastMessageSignatureIn) error
	QueryMessagesInQueue(ctx context.Context, queueTypeName string) ([]chain.MessageWithSignatures, error)
	QueryMessagesForSigning(ctx context.Context, queueTypeName string) ([]chain.QueuedMessage, error)
	QueryGetEVMChainInfos(ctx context.Context) ([]*evmtypes.ChainInfo, error)
	AddMessageEvidence(ctx context.Context, queueTypeName string, messageID uint64, proof []byte) error
	SetPublicAccessData(ctx context.Context, queueTypeName string, messageID uint64, data []byte) error
	QueryGetEVMValsetByID(ctx context.Context, id uint64, chainID string) (*evmtypes.Valset, error)
}

//go:generate mockery --name=EvmFactorier
type EvmFactorier interface {
	Build(
		cfg config.EVM,
		chainReferenceID,
		smartContractID,
		smartContractABIJson,
		smartContractAddress string,
		chainID *big.Int,
		blockHeight int64,
		blockHeightHash common.Hash,
	) (chain.Processor, error)
}

type Relayer struct {
	config config.Root

	palomaClient PalomaClienter

	evmFactory EvmFactorier
}

func New(config config.Root, palomaClient PalomaClienter, evmFactory EvmFactorier) *Relayer {
	return &Relayer{
		config:       config,
		palomaClient: palomaClient,
		evmFactory:   evmFactory,
	}
}

func (r *Relayer) init() error {

	return nil
}
