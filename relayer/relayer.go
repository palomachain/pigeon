package relayer

import (
	"context"

	"github.com/palomachain/pigeon/attest"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/evm"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/config"
	evmtypes "github.com/palomachain/pigeon/types/paloma/x/evm/types"
	valset "github.com/palomachain/pigeon/types/paloma/x/valset/types"
)

type AttestExecutor interface {
	Execute(context.Context, string, attest.Request) (attest.Evidence, error)
}

//go:generate mockery --name=PalomaClienter
type PalomaClienter interface {
	AddExternalChainInfo(ctx context.Context, chainInfos ...paloma.ChainInfoIn) error
	QueryValidatorInfo(ctx context.Context) ([]*valset.ExternalChainInfo, error)
	BroadcastMessageSignatures(ctx context.Context, signatures ...paloma.BroadcastMessageSignatureIn) error
	QueryMessagesInQueue(ctx context.Context, queueTypeName string) ([]chain.MessageWithSignatures, error)
	QueryMessagesForSigning(ctx context.Context, queueTypeName string) ([]chain.QueuedMessage, error)
	QueryGetEVMChainInfos(ctx context.Context) ([]*evmtypes.ChainInfo, error)
}

type Relayer struct {
	config config.Root

	palomaClient PalomaClienter

	attestExecutor AttestExecutor

	evmFactory *evm.Factory
}

func New(config config.Root, palomaClient PalomaClienter, attestExecutor AttestExecutor, evmFactory *evm.Factory) *Relayer {
	return &Relayer{
		config:         config,
		palomaClient:   palomaClient,
		attestExecutor: attestExecutor,
		evmFactory:     evmFactory,
	}
}

func (r *Relayer) init() error {

	return nil
}
