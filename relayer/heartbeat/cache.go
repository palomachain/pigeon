package heartbeat

import (
	"context"
	"sync"
	"time"

	"github.com/palomachain/pigeon/internal/liblog"
	"github.com/sirupsen/logrus"
)

const (
	cMaxCacheRefreshAttempts      int           = 3
	cCacheRefreshIntervalInBlocks int64         = 20
	cDefaultBlockSpeed            time.Duration = time.Millisecond * 1620
)

type keepAliveCache struct {
	retryFalloff        time.Duration
	locker              sync.Locker
	estimatedBlockSpeed time.Duration
	lastBlockHeight     int64
	lastRefresh         time.Time
	lastAliveUntil      int64
	queryBTL            AliveUntilHeightQuery
	queryBH             CurrentHeightQuery
}

func (c *keepAliveCache) get(ctx context.Context) (int64, error) {
	liblog.WithContext(ctx).Infof("%+v", *c)
	liblog.WithContext(ctx).WithFields(logrus.Fields{
		"retryFalloff":        c.retryFalloff,
		"estimatedBlockSpeed": c.estimatedBlockSpeed,
		"lastBlockHeight":     c.lastBlockHeight,
		"lastRefresh":         c.lastRefresh,
		"lastAliveUntil":      c.lastAliveUntil,
	}).Info("cache get")
	if c.isStale() {
		liblog.WithContext(ctx).WithFields(logrus.Fields{
			"retryFalloff":        c.retryFalloff,
			"estimatedBlockSpeed": c.estimatedBlockSpeed,
			"lastBlockHeight":     c.lastBlockHeight,
			"lastRefresh":         c.lastRefresh,
			"lastAliveUntil":      c.lastAliveUntil,
		}).Info("cache is stale")
		err := linearFalloffRetry(ctx, c.locker, "cache refresh", cMaxCacheRefreshAttempts, c.retryFalloff, c.refresh)
		liblog.WithContext(ctx).WithFields(logrus.Fields{
			"retryFalloff":        c.retryFalloff,
			"estimatedBlockSpeed": c.estimatedBlockSpeed,
			"lastBlockHeight":     c.lastBlockHeight,
			"lastRefresh":         c.lastRefresh,
			"lastAliveUntil":      c.lastAliveUntil,
		}).WithError(err).Info("cache refreshed")
		if err != nil {
			return 0, err
		}
	}

	return c.lastAliveUntil, nil
}

func (c *keepAliveCache) refresh(ctx context.Context, locker sync.Locker) error {
	defer locker.Unlock()

	locker.Lock()
	abh, err := c.queryBTL(ctx)
	if err != nil {
		return err
	}

	bh, err := c.queryBH(ctx)
	if err != nil {
		return err
	}

	c.estimatedBlockSpeed = c.estimateBlockSpeed(bh, time.Now().UTC())
	c.lastAliveUntil = abh
	c.lastBlockHeight = bh
	c.lastRefresh = time.Now().UTC()

	return nil
}

func (c *keepAliveCache) isStale() bool {
	if c.estimatedBlockSpeed == 0 || c.lastBlockHeight == 0 || c.lastRefresh.IsZero() {
		return true
	}

	elapsedMs := time.Now().UTC().Sub(c.lastRefresh).Milliseconds()
	estimatedElapsedBlocks := elapsedMs / c.estimatedBlockSpeed.Milliseconds()

	return estimatedElapsedBlocks >= cCacheRefreshIntervalInBlocks
}

func (c *keepAliveCache) estimateBlockSpeed(bh int64, t time.Time) time.Duration {
	if c.lastBlockHeight == 0 || bh == 0 || t.IsZero() {
		// During the first run, we have no historic data to
		// compare to, so we set a rough estimate.
		return cDefaultBlockSpeed
	}

	if t.Before(c.lastRefresh) {
		return cDefaultBlockSpeed
	}

	blockDiff := bh - c.lastBlockHeight
	timeDiff := t.Sub(c.lastRefresh)
	bpms := timeDiff.Milliseconds() / int64(blockDiff)
	return time.Duration(bpms) * time.Millisecond
}

func (c *keepAliveCache) estimateBlockHeight(t time.Time) int64 {
	if c.estimatedBlockSpeed == 0 || c.lastRefresh.IsZero() || t.IsZero() {
		return c.lastBlockHeight
	}

	if t.Before(c.lastRefresh) {
		return c.lastBlockHeight
	}

	timeDiff := t.Sub(c.lastRefresh)
	blockDiff := timeDiff.Milliseconds() / c.estimatedBlockSpeed.Milliseconds()
	return c.lastBlockHeight + blockDiff
}
