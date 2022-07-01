package evm

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/palomachain/pigeon/chain"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

const (
	queueTurnstoneMessage = "evm-turnstone-message"
)

type Processor struct {
	compass          compass
	evmClient        Client
	chainType        string
	chainReferenceID string

	turnstoneEVMContract common.Address
}

func NewProcessor(
	c Client,
	chainReferenceID string,
	smartContractID string,
	smartContractABI abi.ABI,
	smartContractAddress string,
) Processor {
	comp := newCompassClient(
		smartContractAddress,
		smartContractID,
		chainReferenceID,
		smartContractABI,
		c.paloma,
		c,
	)
	return Processor{
		compass:          comp,
		evmClient:        c,
		chainType:        "EVM",
		chainReferenceID: chainReferenceID,
	}
}

var _ chain.Processor = Processor{}

func (p Processor) SupportedQueues() []string {
	return slice.Map(
		[]string{
			queueTurnstoneMessage,
		},
		func(q string) string {
			return fmt.Sprintf("%s:%s:%s", p.chainType, p.chainReferenceID, q)
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
		ChainType:        p.chainType,
		ChainReferenceID: p.chainReferenceID,
		Address:          p.evmClient.addr.Hex(),
		PubKey:           p.evmClient.addr.Bytes(),
	}
}
