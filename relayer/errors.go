package relayer

import (
	"context"
	goerrors "errors"

	"github.com/VolumeFi/whoops"
	"github.com/palomachain/pigeon/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrMissingChainConfig = errors.Unrecoverable(whoops.String("missing chain config"))
	ErrUnknown            = errors.Unrecoverable(whoops.String("unknown errror"))

	ErrInvalidMinOnChainBalance = whoops.Errorf("invalid minOnChainBalance: %s")

	ErrNotAValidatorAccount = whoops.String("not a validator account")

	ErrValidatorIsNotStaking = whoops.String("validator is not staking")
)

func handleProcessError(err error) error {
	switch {
	case err == nil:
		// success
		return nil
	case goerrors.Is(err, context.Canceled):
		log.WithFields(log.Fields{
			"err": err,
		}).Debug("exited from the process loop due the context being canceled")
		return nil
	case goerrors.Is(err, context.DeadlineExceeded):
		log.WithFields(log.Fields{
			"err": err,
		}).Debug("exited from the process loop due the context deadline being exceeded")
		return nil
	case errors.IsUnrecoverable(err):
		// there is no way that we can recover from this
		log.WithFields(log.Fields{
			"err": err,
		}).Error("unrecoverable error returned")
		return err
	default:
		log.WithFields(log.Fields{
			"err": err,
		}).Error("error returned in process loop")
		return nil
	}
}
