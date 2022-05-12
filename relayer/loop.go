package relayer

import (
	"context"
	"fmt"
	"time"

	"github.com/palomachain/sparrow/errors"
	"github.com/vizualni/whoops"
)

const (
	defaultErrorCountToExit = 5

	defaultLoopTimeout = 10 * time.Second
)

// Start starts the relayer. It's responsible for handling the communication
// with Paloma and other chains.
func (r *Relayer) Start(ctx context.Context) error {

	if err := r.init(); err != nil {
		return err
	}

	if err := r.updateValidatorInfo(ctx); err != nil {
		return err
	}

	consecutiveFailures := whoops.Group{}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(defaultLoopTimeout):
			err := r.oneLoopCall(ctx)
			if err == nil {
				// resetting the failures
				if len(consecutiveFailures) > 0 {
					consecutiveFailures = whoops.Group{}
				}
				continue
			}

			if errors.IsUnrecoverable(err) {
				// there is no way that we can recover from this
				return err
			}

			consecutiveFailures.Add(err)

			if len(consecutiveFailures) >= defaultErrorCountToExit {
				return errors.Unrecoverable(consecutiveFailures)
			}
			// TODO: add logger
			fmt.Println("error happened", err)
		}
	}
}

func (r *Relayer) oneLoopCall(ctx context.Context) error {
	var g whoops.Group

	g.Add(r.signMessagesForExecution(ctx,
		consensusExecuteSmartContract,
		consensusUpdateValset,
	))
	g.Add(r.queryConcencusReachedMessages(ctx))

	if g.Err() {
		return g
	}

	return nil
}
