package relayer

import (
	"context"

	"github.com/palomachain/pigeon/attest"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/chain/paloma"
	"github.com/palomachain/pigeon/config"
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
}

type Relayer struct {
	config config.Root

	palomaClient PalomaClienter

	attestExecutor AttestExecutor

	processors map[string]chain.Processor
}

func New(config config.Root, palomaClient PalomaClienter, attestExecutor AttestExecutor, processors map[string]chain.Processor) *Relayer {
	return &Relayer{
		config:         config,
		palomaClient:   palomaClient,
		attestExecutor: attestExecutor,
		processors:     processors,
	}
}

func (r *Relayer) init() error {

	return nil
}
