package chain

import "context"

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

type Relayer interface {
	SupportedQueues() []string
	UpdateValset(context.Context, ValsetUpdate) (ValsetUpdateResponse, error)
	ExecuteArbitraryMessage(context.Context, ArbitraryMessage) (ArbitraryMessageResponse, error)
}
