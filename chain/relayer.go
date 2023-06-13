package chain

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/palomachain/pigeon/health"
)

type QueuedMessage struct {
	ID               uint64
	Nonce            []byte
	BytesToSign      []byte
	PublicAccessData []byte
	ErrorData        []byte
	Msg              any
}

type SignedQueuedMessage struct {
	QueuedMessage
	Signature       []byte
	SignedByAddress string
}

type MessageToProcess struct {
	QueueTypeName string
	Msg           QueuedMessage
}

type ValidatorSignature struct {
	ValAddress      sdk.ValAddress
	Signature       []byte
	SignedByAddress string
	PublicKey       []byte
}

type MessageWithSignatures struct {
	QueuedMessage
	Signatures []ValidatorSignature
}

type ExternalAccount struct {
	ChainType        string
	ChainReferenceID string

	Address string
	PubKey  []byte
}

type ChainInfo interface {
	ChainReferenceID() string
	ChainID() string
	ChainType() string
}

//go:generate mockery --name=Processor
type Processor interface {
	health.Checker
	// SupportedQueues is a list of consensus queues that this processor supports and expects to work with.
	SupportedQueues() []string

	ExternalAccount() ExternalAccount

	// SignMessages takes a list of messages and signs them via their key.
	SignMessages(ctx context.Context, queueTypeName string, messages ...QueuedMessage) ([]SignedQueuedMessage, error)

	// ProcessMessages will receive messages from the current queues and it's on the implementation
	// to ensure that there are enough signatures for consensus.
	ProcessMessages(context.Context, string, []MessageWithSignatures) error

	// ProvideEvidence takes a queue name and a list of messages that have already been executed. This
	// takes the "public evidence" from the message and gets the information back to the Paloma.
	ProvideEvidence(context.Context, string, []MessageWithSignatures) error

	// it verifies if it's being connected to the right chain
	IsRightChain(ctx context.Context) error
}

type ProcessorBuilder interface {
	Build(ChainInfo) (Processor, error)
}
