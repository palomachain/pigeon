package relayer

import (
	"context"
	"strings"
	"sync"

	"github.com/palomachain/pigeon/internal/liblog"
	log "github.com/sirupsen/logrus"
)

func (r *Relayer) keepAlive(ctx context.Context, locker sync.Locker) error {
	logger := liblog.WithContext(ctx)
	logger.Debug("querying get alive time")

	aliveUntil, err := r.palomaClient.QueryGetValidatorAliveUntilBlockHeight(ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "validator is not in keep alive store") {
			log.WithError(err).Error("error while getting the alive time for a validator")
			return err
		}
	}

	bh, err := r.palomaClient.BlockHeight(ctx)
	if err != nil {
		return err
	}

	btl := aliveUntil - bh
	sendKeepAlive := btl < r.relayerConfig.KeepAliveBlockThreshold
	log.WithFields(log.Fields{
		"alive-until-bh":         aliveUntil,
		"current-bh":             bh,
		"btl":                    btl,
		"should-send-keep-alive": sendKeepAlive,
	}).Debug("checking keep alive")

	if sendKeepAlive {
		locker.Lock()
		err := r.palomaClient.KeepValidatorAlive(ctx, r.appVersion)
		locker.Unlock()
		if err != nil {
			log.WithError(err).Error("error while trying to keep pigeon alive")
			return err
		}
	}

	return nil
}
