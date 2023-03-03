package relayer

import (
	"github.com/VolumeFi/whoops"
	"github.com/palomachain/pigeon/errors"
)

var (
	ErrMissingChainConfig = errors.Unrecoverable(whoops.String("missing chain config"))
	ErrUnknown            = errors.Unrecoverable(whoops.String("unknown errror"))

	ErrInvalidMinOnChainBalance = whoops.Errorf("invalid minOnChainBalance: %s")

	ErrNotAValidatorAccount = whoops.String("not a validator account")

	ErrValidatorIsNotStaking = whoops.String("validator is not staking")
)
