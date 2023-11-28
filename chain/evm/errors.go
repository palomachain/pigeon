package evm

import (
	"github.com/VolumeFi/whoops"
)

const (
	ErrSmartContractNotFound     = whoops.Errorf("smart contract %s was not found")
	ErrInvalidAddress            = whoops.Errorf("provided address: '%s' is not valid")
	ErrAddressNotFoundInKeyStore = whoops.Errorf("address: '%s' not found in keystore: %s")
	ErrUnsupportedMessageType    = whoops.Errorf("unsupported message type: %T")
	ErrABINotInitialized         = whoops.String("ABI is not initialized")

	ErrEvm = whoops.String("EVM related error")

	ErrNoConsensus         = whoops.String("no consensus reached")
	ErrCallAlreadyExecuted = whoops.String("message with identical ID already on target chain")

	ErrCouldntFindBlockWithTime = whoops.String("couldn't find block")
)

var ErrStartingBlockIsInTheFuture = ErrCouldntFindBlockWithTime.WrapS("starting height's block time is set in future")

const (
	FieldMessageID   whoops.Field[uint64] = "message id"
	FieldMessageType whoops.Field[any]    = "message type"
)
