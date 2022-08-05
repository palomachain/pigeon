package relayer

import (
	"github.com/palomachain/pigeon/errors"
	"github.com/vizualni/whoops"
)

var (
	ErrMissingChainConfig = errors.Unrecoverable(whoops.String("missing chain config"))
	ErrUnknown            = errors.Unrecoverable(whoops.String("unknown errror"))
)
