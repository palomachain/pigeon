package health

import (
	"context"
	"time"

	"github.com/palomachain/pigeon/util/channels"
	log "github.com/sirupsen/logrus"
)

type PalomaStatuser interface {
	PalomaStatus(ctx context.Context) error
}

func WaitForPaloma(ctx context.Context, ps PalomaStatuser) error {
	const checkTimeout = 5 * time.Second
	t := time.NewTicker(checkTimeout)
	t1 := make(chan time.Time, 1)
	t1 <- time.Time{}
	ch := channels.FanIn(t.C, t1)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
			statusErr := ps.PalomaStatus(ctx)
			if statusErr == nil {
				// good!
				// paloma was detected to be running
				// giving some time for paloma to start
				log.Info("paloma detected being online. Sleeping for 5 seconds to give it time to properly start.")
				time.Sleep(5 * time.Second)
				return nil
			}
			log.WithError(statusErr).Error("waiting for paloma to start running")
		}
	}
}

func CancelContextIfPalomaIsDown(ctx context.Context, ps PalomaStatuser) context.Context {
	ctx, cancelFnc := context.WithCancel(ctx)

	go func() {
		defer cancelFnc()
		const checkTimeout = 5 * time.Second
		t := time.NewTicker(checkTimeout)
		t1 := make(chan time.Time, 1)
		t1 <- time.Time{}
		ch := channels.FanIn(t.C, t1)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ch:
				statusErr := ps.PalomaStatus(ctx)
				if statusErr != nil {
					log.WithError(statusErr).Error("unable to get paloma status")
					return
				}
			}
		}
	}()

	return ctx
}
