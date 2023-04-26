package relayer

import (
	"context"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	proto "github.com/cosmos/gogoproto/proto"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/config"
	evmtypes "github.com/palomachain/pigeon/types/paloma/x/evm/types"
	valset "github.com/palomachain/pigeon/types/paloma/x/valset/types"
	utiltime "github.com/palomachain/pigeon/util/time"
)

//go:generate mockery --name=PalomaClienter
type PalomaClienter interface {
	AddExternalChainInfo(ctx context.Context, chainInfos ...paloma.ChainInfoIn) error
	QueryValidatorInfo(ctx context.Context) ([]*valset.ExternalChainInfo, error)
	BroadcastMessageSignatures(ctx context.Context, signatures ...paloma.BroadcastMessageSignatureIn) error
	QueryMessagesInQueue(ctx context.Context, queueTypeName string) ([]chain.MessageWithSignatures, error)
	QueryMessagesForSigning(ctx context.Context, queueTypeName string) ([]chain.QueuedMessage, error)
	QueryGetEVMChainInfos(ctx context.Context) ([]*evmtypes.ChainInfo, error)
	AddMessageEvidence(ctx context.Context, queueTypeName string, messageID uint64, proof proto.Message) error
	SetPublicAccessData(ctx context.Context, queueTypeName string, messageID uint64, data []byte) error
	QueryGetEVMValsetByID(ctx context.Context, id uint64, chainID string) (*evmtypes.Valset, error)
	GetValidatorAddress() sdk.ValAddress
	GetValidator(ctx context.Context) (*stakingtypes.Validator, error)

	BlockHeight(context.Context) (int64, error)
	QueryGetSnapshotByID(ctx context.Context, id uint64) (*valset.Snapshot, error)

	QueryGetValidatorAliveUntil(ctx context.Context) (time.Time, error)
	KeepValidatorAlive(ctx context.Context) error
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
		minOnChainBalance *big.Int,
	) (chain.Processor, error)
}

type Relayer struct {
	config config.Root

	palomaClient PalomaClienter

	evmFactory EvmFactorier

	relayerConfig Config

	time utiltime.Time
}

type Config struct {
	KeepAliveLoopTimeout time.Duration
	KeepAliveThreshold   time.Duration
}

func New(config config.Root, palomaClient PalomaClienter, evmFactory EvmFactorier, customTime utiltime.Time, cfg Config) *Relayer {
	return &Relayer{
		config:        config,
		palomaClient:  palomaClient,
		evmFactory:    evmFactory,
		time:          customTime,
		relayerConfig: cfg,
	}
}
