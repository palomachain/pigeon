package chain

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ValsetUpdate struct{}

type ValsetUpdateResponse struct{}

type ArbitraryMessage struct {
	QueueName string
}

type ArbitraryMessageResponse struct{}

// TODO: thing such as registering itself with the paloma
type PalomaCompatibler interface{}

type CustomSigner interface {
	// Sign signs a message in a way that can be understood by the target
	// chain.
	Sign(ctx context.Context, msg []byte) (sig []byte, err error)
}

type QueuedMessage struct {
	ID          uint64
	Nonce       []byte
	BytesToSign []byte
	Msg         any
}

type MessageToProcess struct {
	QueueTypeName string
	Msg           QueuedMessage
}

type ValidatorSignature struct {
	ValAddress sdk.ValAddress
	Signature  []byte
}
type ConsensusReachedMsg struct {
	ID         string
	Nonce      []byte
	Signatures []ValidatorSignature
	Msg        any
}

type Processor interface {
	SupportedQueues() []string
	SignMessages(ctx context.Context, queueTypeName string, messages ...QueuedMessage)
	RelayMessages(ctx context.Context, queueTypeName string, messages ...ConsensusReachedMsg)
}

type Relayer interface {
	UpdateValset(context.Context, ValsetUpdate) (ValsetUpdateResponse, error)
	ExecuteArbitraryMessage(context.Context, ArbitraryMessage) (ArbitraryMessageResponse, error)
}
