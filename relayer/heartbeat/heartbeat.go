package heartbeat

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/palomachain/pigeon/internal/liblog"
	log "github.com/sirupsen/logrus"
)

const (
	cMaxSendKeepAliveAttempts int           = 3
	cDefaultRetryFalloff      time.Duration = time.Second
)

type (
	AliveUntilHeightQuery func(context.Context) (int64, error)
	CurrentHeightQuery    func(context.Context) (int64, error)
	KeepAliveCall         func(context.Context, string) error
)

type Heart struct {
	sendKeepAlive      KeepAliveCall
	keepAliveThreshold int64
	appVersion         string
	c                  *keepAliveCache
	retryFalloff       time.Duration
}

func New(i AliveUntilHeightQuery, j CurrentHeightQuery, k KeepAliveCall, keepAliveThreshold int64, appVersion string, locker sync.Locker) *Heart {
	cache := &keepAliveCache{
		locker:       locker,
		queryBTL:     i,
		queryBH:      j,
		retryFalloff: cDefaultRetryFalloff,
	}
	return &Heart{
		sendKeepAlive:      k,
		keepAliveThreshold: keepAliveThreshold,
		appVersion:         appVersion,
		c:                  cache,
		retryFalloff:       cDefaultRetryFalloff,
	}
}

func (m *Heart) Beat(ctx context.Context, locker sync.Locker) error {
	logger := liblog.WithContext(ctx)

	aliveUntil, err := m.c.get(ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "validator is not in keep alive store") {
			logger.WithError(err).Error("error while getting the alive time for a validator")
			return err
		}
	}

	bh := m.c.estimateBlockHeight(time.Now().UTC())
	btl := aliveUntil - bh
	sendKeepAlive := btl < m.keepAliveThreshold
	logger.WithFields(log.Fields{
		"alive-until-bh":         aliveUntil,
		"current-bh":             bh,
		"btl":                    btl,
		"should-send-keep-alive": sendKeepAlive,
	}).Debug("checking keep alive")

	if sendKeepAlive {
		return linearFalloffRetry(ctx, locker, "keep alive", cMaxSendKeepAliveAttempts, cDefaultRetryFalloff, m.trySendKeepAlive)
	}

	return nil
}

func (m *Heart) SetRetryFalloff(falloff time.Duration) {
	m.retryFalloff = falloff
	m.c.retryFalloff = falloff
}

// Will be blocking during retries. Make sure to always call in Goroutine.
func (m *Heart) trySendKeepAlive(ctx context.Context, locker sync.Locker) (err error) {
	logger := liblog.WithContext(ctx)

	locker.Lock()
	err = m.sendKeepAlive(ctx, m.appVersion)
	locker.Unlock()
	if err == nil {
		return nil
	}

	logger.WithError(err).Error("Error while trying to keep pigeon alive.")
	return err
}

func linearFalloffRetry(ctx context.Context, locker sync.Locker, name string, maxRetries int, baseFalloff time.Duration, fn func(context.Context, sync.Locker) error) (err error) {
	logger := liblog.WithContext(ctx)
	retries := 0
	falloff := baseFalloff

	for retries < maxRetries {
		err = fn(ctx, locker)
		if err == nil {
			return nil
		}

		logger.WithError(err).Warnf(
			"%s: Attempt [%d/%d] failed to send keep alive call. Will retry in %v",
			name,
			retries+1,
			cMaxSendKeepAliveAttempts,
			falloff)

		time.Sleep(falloff)

		retries++
		falloff += falloff
	}

	return err
}
