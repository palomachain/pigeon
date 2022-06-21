package evm

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/palomachain/sparrow/chain"
	"github.com/palomachain/sparrow/util/slice"
	log "github.com/sirupsen/logrus"
)

const (
	queueArbitraryLogic   = "evm-arbitrary-smart-contract-call"
	queueTurnstoneMessage = "evm-turnstone-message"
)

type Processor struct {
	compass   compass
	evmClient Client
	chainType string
	chainID   string

	turnstoneEVMContract common.Address
}

func NewProcessor(c Client, chainID string) Processor {
	comp := newCompassClient(
		c.config.SmartContractAddress,
		c.config.CompassID,
		c.internalChainID,
		c.smartContractAbi,
		c.paloma,
		c,
	)
	return Processor{
		compass:   comp,
		evmClient: c,
		chainType: "EVM",
		chainID:   chainID,
	}
}

var _ chain.Processor = Processor{}

func (p Processor) SupportedQueues() []string {
	return slice.Map(
		[]string{
			// queueArbitraryLogic,
			queueTurnstoneMessage,
		},
		func(q string) string {
			return fmt.Sprintf("%s:%s:%s", p.chainType, p.chainID, q)
		},
	)
}

func (p Processor) SignMessages(ctx context.Context, queueTypeName string, messages ...chain.QueuedMessage) ([]chain.SignedQueuedMessage, error) {
	return slice.MapErr(messages, func(msg chain.QueuedMessage) (chain.SignedQueuedMessage, error) {
		msgBytes := crypto.Keccak256(
			append(
				[]byte(SignedMessagePrefix),
				msg.BytesToSign...,
			),
		)
		sig, err := p.evmClient.sign(ctx, msgBytes)
		log.WithFields(log.Fields{
			"msg": msg,
			"sig": sig,
			"err": err,
		}).Info("signing a message")

		if err != nil {
			return chain.SignedQueuedMessage{}, err
		}

		return chain.SignedQueuedMessage{
			QueuedMessage:   msg,
			Signature:       sig,
			SignedByAddress: p.evmClient.addr.Hex(),
		}, nil
	},
	)

}

func (p Processor) ProcessMessages(ctx context.Context, queueTypeName string, msgs []chain.MessageWithSignatures) error {
	switch {
	case strings.HasSuffix(queueTypeName, queueTurnstoneMessage):
		return p.compass.processMessages(
			ctx,
			queueTypeName,
			msgs,
		)
	default:
		return chain.ErrProcessorDoesNotSupportThisQueue.Format(queueTypeName)
	}
}

func (p Processor) ExternalAccount() chain.ExternalAccount {
	return chain.ExternalAccount{
		ChainType: p.chainType,
		ChainID:   p.chainID,
		Address:   p.evmClient.addr.Hex(),
		PubKey:    p.evmClient.addr.Bytes(),
	}
}
