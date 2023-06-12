package relayer

import (
	"context"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

func (r *Relayer) keepAlive(ctx context.Context, locker sync.Locker) error {
	log.Debug("querying get alive time")
	aliveUntil, err := r.palomaClient.QueryGetValidatorAliveUntil(ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "validator is not in keep alive store") {
			log.WithError(err).Error("error while getting the alive time for a validator")
			return err
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
			return err
		}
	}
	return nil
}
