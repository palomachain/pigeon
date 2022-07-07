package evm

import (
	"github.com/vizualni/whoops"
)

const (
	ErrSmartContractNotFound     = whoops.Errorf("smart contract %s was not found")
	ErrInvalidAddress            = whoops.Errorf("provided address: '%s' is not valid")
	ErrAddressNotFoundInKeyStore = whoops.Errorf("address: '%s' not found in keystore: %s")
	ErrUnsupportedMessageType    = whoops.Errorf("unsupported message type: %T")
	ErrABINotInitialized         = whoops.String("ABI is not initialized")

	ErrNoConsensus = whoops.String("no consensus reached")
)
