package relayer

import (
	"context"
	"fmt"
	"time"

	"github.com/vizualni/whoops"
	"github.com/palomachain/sparrow/errors"
)

const (
	// TODO: FILL IN THE REAL QUEUE TYPE NAMES
	queueTypeNameMessageExecution = "a"
	queueTypeNameValsetUpdate     = "b"

	defaultErrorCountToExit = 5

	defaultLoopTimeout = 10 * time.Second
)

// Start starts the relayer. It's responsible for handling the communication
// with Paloma and other chains.
func (r *Relayer) Start(ctx context.Context) error {
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

	g.Add(r.signMessagesForExecution(ctx, queueTypeNameMessageExecution, queueTypeNameValsetUpdate))
	g.Add(r.queryConcencusReachedMessages(ctx))

	if g.Err() {
		return g
	}

	return nil
}
