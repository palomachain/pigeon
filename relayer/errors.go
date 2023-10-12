package relayer

import (
	"context"
	goerrors "errors"

	"github.com/VolumeFi/whoops"
	"github.com/palomachain/pigeon/errors"
	"github.com/palomachain/pigeon/internal/liblog"
)

var (
	ErrMissingChainConfig = errors.Unrecoverable(whoops.String("missing chain config"))
	ErrUnknown            = errors.Unrecoverable(whoops.String("unknown errror"))

	ErrInvalidMinOnChainBalance = whoops.Errorf("invalid minOnChainBalance: %s")

	ErrNotAValidatorAccount = whoops.String("not a validator account")

	ErrValidatorIsNotStaking = whoops.String("validator is not staking")
)

func handleProcessError(ctx context.Context, err error) error {
	if err == nil {
		return err
	}

	logger := liblog.WithContext(ctx)
	if goerrors.Is(err, context.Canceled) || goerrors.Is(err, context.DeadlineExceeded) {
		logger.WithError(err).Warn("exicted from the loop because of context error")
		return nil
	}

	if errors.IsUnrecoverable(err) {
		logger.WithError(err).Error("unrecoverable error returned")
		return err
	}

	logger.WithError(err).Error("error returned in process loop")
	return nil
}
