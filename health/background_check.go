package health

import (
	"context"
	"time"

	"github.com/VolumeFi/whoops"
	"github.com/palomachain/pigeon/util/channels"
	log "github.com/sirupsen/logrus"
)

const (
	checkTimeout = time.Minute
)

type Errors struct {
	All whoops.Group
}

func (e *Errors) Error() string {
	return e.All.Error()
}

type Service struct {
	Checks []Checker
}

func (s Service) HealthCheckInBackground(
	ctx context.Context,
) {
	t := time.NewTicker(checkTimeout)
	t1 := make(chan time.Time, 1)
	t1 <- time.Time{}
	ch := channels.FanIn(t.C, t1)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ch:
			s.Check(ctx)
		}
	}
}

func (s Service) Check(ctx context.Context) {
	var g whoops.Group

	for _, hc := range s.Checks {
		g.Add(hc.HealthCheck(ctx))
	}

	if g.Err() {
		for _, err := range g {
			log.WithError(err).Error("health check failed")
		}
		log.Fatal("exiting due to health check failures")
	}
}

func (s Service) BootChecker(ctx context.Context) {
	var g whoops.Group

	for _, hc := range s.Checks {
		b, ok := hc.(BootChecker)
		if !ok {
			continue
		}
		g.Add(b.BootHealthCheck(ctx))
	}

	if g.Err() {
		for _, err := range g {
			log.WithError(err).Warn("boot health check failed")
		}
	}
}
