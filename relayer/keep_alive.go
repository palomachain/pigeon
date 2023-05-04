package relayer

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/palomachain/pigeon/util/channels"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) startKeepAlive(ctx context.Context, locker *sync.Mutex) {
	log.Debug("starting keep alive loop")
	defer func() {
		log.Debug("existing keep alive loop")
	}()
	ticker := time.NewTicker(r.relayerConfig.KeepAliveLoopTimeout)
	defer ticker.Stop()

	checkNow := make(chan time.Time, 1)
	checkNow <- time.Time{}
	tickerCh := channels.FanIn(ticker.C, checkNow)
	defer func() {
		log.Info("exiting keep alive loop")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case _, chOpen := <-tickerCh:
			if !chOpen {
				return
			}
			if ctx.Err() != nil {
				return
			}
			log.Debug("querying get alive time")
			aliveUntil, err := r.palomaClient.QueryGetValidatorAliveUntil(ctx)
			if err != nil {
				if !strings.Contains(err.Error(), "validator is not in keep alive store") {
					log.WithError(err).Error("error while getting the alive time for a validator")
					continue
				}
			}
			now := r.time.Now().UTC()
			ttl := aliveUntil.Sub(now)
			sendKeepAlive := ttl < r.relayerConfig.KeepAliveThreshold
			log.WithFields(log.Fields{
				"alive-until":            aliveUntil,
				"time-now":               now,
				"ttl":                    ttl,
				"should-send-keep-alive": sendKeepAlive,
			}).Debug("checking keep alive")
			if sendKeepAlive {
				locker.Lock()
				err := r.palomaClient.KeepValidatorAlive(ctx)
				locker.Unlock()
				if err != nil {
					log.WithError(err).Error("error while trying to keep pigeon alive")
					continue
				}
			}
		}
	}
}
