package chain

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gravity "github.com/palomachain/paloma/x/gravity/types"

	"github.com/palomachain/pigeon/health"
	"github.com/palomachain/pigeon/internal/queue"
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

type SignedGravityOutgoingTxBatch struct {
	gravity.OutgoingTxBatch
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

type SignedEntity interface {
	GetSignatures() []ValidatorSignature
	GetBytes() []byte
}

type MessageWithSignatures struct {
	QueuedMessage
	Signatures []ValidatorSignature
}

func (msg MessageWithSignatures) GetSignatures() []ValidatorSignature {
	return msg.Signatures
}

func (msg MessageWithSignatures) GetBytes() []byte {
	return msg.BytesToSign
}

type GravityBatchWithSignatures struct {
	gravity.OutgoingTxBatch
	Signatures []ValidatorSignature
}

func (gb GravityBatchWithSignatures) GetSignatures() []ValidatorSignature {
	return gb.Signatures
}

func (gb GravityBatchWithSignatures) GetBytes() []byte {
	return gb.GetBytesToSign()
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
	// GetChainReferenceID returns the chain reference ID against which the processor is running.
	GetChainReferenceID() string

	// SupportedQueues is a list of consensus queues that this processor supports and expects to work with.
	SupportedQueues() []string

	ExternalAccount() ExternalAccount

	// SignMessages takes a list of messages and signs them via their key.
	SignMessages(ctx context.Context, messages ...QueuedMessage) ([]SignedQueuedMessage, error)

	// ProcessMessages will receive messages from the current queues and it's on the implementation
	// to ensure that there are enough signatures for consensus.
	ProcessMessages(context.Context, queue.TypeName, []MessageWithSignatures) error

	// ProvideEvidence takes a queue name and a list of messages that have already been executed. This
	// takes the "public evidence" from the message and gets the information back to the Paloma.
	ProvideEvidence(context.Context, queue.TypeName, []MessageWithSignatures) error

	// it verifies if it's being connected to the right chain
	IsRightChain(ctx context.Context) error

	GravitySignBatches(ctx context.Context, batches ...gravity.OutgoingTxBatch) ([]SignedGravityOutgoingTxBatch, error)
	GravityRelayBatches(ctx context.Context, batches []GravityBatchWithSignatures) error
}

type ProcessorBuilder interface {
	Build(ChainInfo) (Processor, error)
}
