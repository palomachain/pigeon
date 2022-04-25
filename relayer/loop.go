package relayer

import (
	"context"
	"time"
)

const (
	queueTypeNameMessageExecution = "a"
	queueTypeNameValsetUpdate     = "b"
)

func (r relayer) mainLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(999 * time.Second):
			r.loopLogic(ctx)
		}
	}
}

func (r relayer) loopLogic(ctx context.Context) {
	err := r.signMessagesForExecution(ctx, queueTypeNameMessageExecution, queueTypeNameValsetUpdate)
	if err != nil {
		panic(err)
	}
}
