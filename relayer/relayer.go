package relayer

import (
	"context"
	gravity "github.com/palomachain/paloma/x/gravity/types"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	proto "github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/palomachain/paloma/x/evm/types"
	valset "github.com/palomachain/paloma/x/valset/types"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/config"
	"github.com/palomachain/pigeon/internal/mev"
	utiltime "github.com/palomachain/pigeon/util/time"
)

//go:generate mockery --name=PalomaClienter
type PalomaClienter interface {
	AddExternalChainInfo(ctx context.Context, chainInfos ...paloma.ChainInfoIn) error
	QueryValidatorInfo(ctx context.Context) ([]*valset.ExternalChainInfo, error)
	BroadcastMessageSignatures(ctx context.Context, signatures ...paloma.BroadcastMessageSignatureIn) error
	QueryMessagesForAttesting(ctx context.Context, queueTypeName string) ([]chain.MessageWithSignatures, error)
	QueryMessagesForRelaying(ctx context.Context, queueTypeName string) ([]chain.MessageWithSignatures, error)
	QueryMessagesForSigning(ctx context.Context, queueTypeName string) ([]chain.QueuedMessage, error)
	QueryGetEVMChainInfos(ctx context.Context) ([]*evmtypes.ChainInfo, error)
	AddMessageEvidence(ctx context.Context, queueTypeName string, messageID uint64, proof proto.Message) error
	SetPublicAccessData(ctx context.Context, queueTypeName string, messageID uint64, data []byte) error
	SetErrorData(ctx context.Context, queueTypeName string, messageID uint64, data []byte) error
	QueryGetEVMValsetByID(ctx context.Context, id uint64, chainID string) (*evmtypes.Valset, error)
	GetValidatorAddress() sdk.ValAddress
	GetValidator(ctx context.Context) (*stakingtypes.Validator, error)
	GetCreator() string
	BlockHeight(context.Context) (int64, error)
	QueryGetSnapshotByID(ctx context.Context, id uint64) (*valset.Snapshot, error)

	QueryGetValidatorAliveUntil(ctx context.Context) (time.Time, error)
	KeepValidatorAlive(ctx context.Context, appVersion string) error

	GravityRequestBatch(ctx context.Context, chainReferenceId string) error
	GravityQueryLastUnsignedBatch(ctx context.Context, chainReferenceID string) ([]gravity.OutgoingTxBatch, error)
	GravityConfirmBatches(ctx context.Context, signatures ...chain.SignedGravityOutgoingTxBatch) error
	GravityQueryBatchesForRelaying(ctx context.Context, chainReferenceID string) ([]chain.GravityBatchWithSignatures, error)
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
		mevClient mev.Client,
	) (chain.Processor, error)
}

type Relayer struct {
	config config.Root

	palomaClient PalomaClienter

	evmFactory EvmFactorier
	mevClient  mev.Client

	relayerConfig Config

	time utiltime.Time

	chainsInfos []evmtypes.ChainInfo
	processors  []chain.Processor

	staking bool

	appVersion string
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
		staking:       false,
	}
}

func (r *Relayer) SetAppVersion(appVersion string) {
	r.appVersion = appVersion
}

func (r *Relayer) SetMevClient(c mev.Client) {
	r.mevClient = c
}
