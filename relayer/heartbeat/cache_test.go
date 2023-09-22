package heartbeat

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	tm := time.Now().UTC()

	t.Run("estimateBlockHeight", func(t *testing.T) {
		t.Run("with empty block speed", func(t *testing.T) {
			c := &keepAliveCache{
				lastRefresh: tm,
			}
			i := c.estimateBlockHeight(tm)
			require.Equal(t, int64(0), i, "must return last blockheight")
		})

		t.Run("with empty last refresh timestamp", func(t *testing.T) {
			c := &keepAliveCache{
				estimatedBlockSpeed: time.Second,
			}
			i := c.estimateBlockHeight(tm)
			require.Equal(t, c.lastBlockHeight, i, "must return last blockheight")
		})

		t.Run("with empty reference time argument", func(t *testing.T) {
			c := &keepAliveCache{
				lastBlockHeight: 10,
				lastRefresh:     tm,
			}
			i := c.estimateBlockHeight(time.Time{})
			require.Equal(t, c.lastBlockHeight, i, "must return last blockheight")
		})

		t.Run("with reference time before last refresh time", func(t *testing.T) {
			c := &keepAliveCache{
				lastBlockHeight:     10,
				estimatedBlockSpeed: time.Millisecond * 1620,
				lastRefresh:         tm,
			}
			i := c.estimateBlockHeight(tm.Add(time.Hour * -1))
			require.Equal(t, c.lastBlockHeight, i, "must return last blockheight")
		})

		t.Run("with valid parameters", func(t *testing.T) {
			c := &keepAliveCache{
				lastBlockHeight:     10,
				estimatedBlockSpeed: time.Millisecond * 1620,
				lastRefresh:         tm.Add(time.Minute * -1),
			}
			i := c.estimateBlockHeight(tm)
			bh := int64(37 + 10) // 1.62s per block = 37 per minute + 10 for last height
			require.Equal(t, bh, i, "must return the correct calculated block height")
		})
	})

	t.Run("estimateBlockSpeed", func(t *testing.T) {
		t.Run("with empty last block height", func(t *testing.T) {
			c := &keepAliveCache{}
			i := c.estimateBlockSpeed(10, tm)
			require.Equal(t, cDefaultBlockSpeed, i, "must return default block speed")
		})

		t.Run("with empty reference time argument", func(t *testing.T) {
			c := &keepAliveCache{
				lastBlockHeight: 10,
			}
			i := c.estimateBlockSpeed(10, time.Time{})
			require.Equal(t, cDefaultBlockSpeed, i, "must return default block speed")
		})

		t.Run("with zero block height difference", func(t *testing.T) {
			c := &keepAliveCache{
				lastBlockHeight: 10,
			}
			i := c.estimateBlockSpeed(0, tm)
			require.Equal(t, cDefaultBlockSpeed, i, "must return default block speed")
		})

		t.Run("with reference time before last refresh time", func(t *testing.T) {
			c := &keepAliveCache{
				lastBlockHeight: 10,
				lastRefresh:     tm,
			}
			i := c.estimateBlockSpeed(10, tm.Add(time.Hour*-1))
			require.Equal(t, cDefaultBlockSpeed, i, "must return default block speed")
		})

		t.Run("with unrealistically low block height difference", func(t *testing.T) {
			c := &keepAliveCache{
				lastBlockHeight: 10,
				lastRefresh:     tm.Add(time.Minute * -1),
			}
			i := c.estimateBlockSpeed(29, tm)
			require.Equal(t, cDefaultBlockSpeed, i, "must return default block speed")
		})

		t.Run("with unrealistically high time difference", func(t *testing.T) {
			c := &keepAliveCache{
				lastBlockHeight: 10,
				lastRefresh:     tm.Add(time.Hour * -1),
			}
			i := c.estimateBlockSpeed(60, tm)
			require.Equal(t, cDefaultBlockSpeed, i, "must return default block speed")
		})

		t.Run("with valid parameters", func(t *testing.T) {
			c := &keepAliveCache{
				lastBlockHeight: 10,
				lastRefresh:     tm.Add(time.Minute * -1),
			}
			i := c.estimateBlockSpeed(47, tm)
			bs := time.Millisecond * 1621
			require.Equal(t, bs, i, "must return the correct calculated block speed")
		})
	})

	t.Run("isStale", func(t *testing.T) {
		t.Run("with missing invalidated flag set", func(t *testing.T) {
			c := &keepAliveCache{
				lastBlockHeight:     10,
				lastRefresh:         tm,
				estimatedBlockSpeed: time.Second,
			}
			c.invalidate()
			require.True(t, c.isStale(), "must return true")
		})

		t.Run("with missing estimated block speed", func(t *testing.T) {
			c := &keepAliveCache{
				lastBlockHeight: 10,
				lastRefresh:     tm,
			}
			require.True(t, c.isStale(), "must return true")
		})

		t.Run("with missing last refresh time", func(t *testing.T) {
			c := &keepAliveCache{
				lastBlockHeight:     10,
				estimatedBlockSpeed: time.Second,
			}
			require.True(t, c.isStale(), "must return true")
		})

		t.Run("with missing last block height", func(t *testing.T) {
			c := &keepAliveCache{
				lastRefresh:         tm,
				estimatedBlockSpeed: time.Second,
			}
			require.True(t, c.isStale(), "must return true")
		})

		t.Run("with estimated elapsed blocks less than cache refresh interval", func(t *testing.T) {
			c := &keepAliveCache{
				lastBlockHeight:     10,
				estimatedBlockSpeed: time.Second,
				lastRefresh:         tm.Add(time.Second * -10),
			}
			require.False(t, c.isStale(), "must return false")
		})

		t.Run("with estimated elapsed blocks greater than or equal to cache refresh interval", func(t *testing.T) {
			c := &keepAliveCache{
				lastBlockHeight:     10,
				estimatedBlockSpeed: time.Second,
				lastRefresh:         tm.Add(time.Second * -20),
			}
			require.True(t, c.isStale(), "must return true")
		})
	})

	t.Run("refresh", func(t *testing.T) {
		errTest := fmt.Errorf("fail")
		ctx := context.Background()

		t.Run("with error from alive until height query", func(t *testing.T) {
			c := &keepAliveCache{
				queryBTL: func(ctx context.Context) (int64, error) {
					return 0, errTest
				},
			}
			err := c.refresh(ctx, &sync.Mutex{})
			require.Error(t, err, "must return error")
		})

		t.Run("with error from current height query", func(t *testing.T) {
			c := &keepAliveCache{
				queryBTL: func(ctx context.Context) (int64, error) {
					return 10, nil
				},
				queryBH: func(ctx context.Context) (int64, error) {
					return 0, errTest
				},
			}
			err := c.refresh(ctx, &sync.Mutex{})
			require.Error(t, err, "must return error")
		})

		t.Run("with valid parameters", func(t *testing.T) {
			c := &keepAliveCache{
				queryBTL: func(ctx context.Context) (int64, error) {
					return 47, nil
				},
				queryBH: func(ctx context.Context) (int64, error) {
					return 10, nil
				},
			}
			err := c.refresh(ctx, &sync.Mutex{})
			require.NoError(t, err, "must not return error")
			require.Equal(t, int64(47), c.lastAliveUntil, "must set last alive until")
			require.Equal(t, int64(10), c.lastBlockHeight, "must set last block height")
			require.LessOrEqual(t, c.lastRefresh.Sub(tm), time.Duration(time.Second), "must set last refresh")
		})
	})

	t.Run("get", func(t *testing.T) {
		ctx := context.Background()

		t.Run("with stale cache", func(t *testing.T) {
			c := &keepAliveCache{
				lastRefresh: tm.Add(time.Second * -20),
				locker:      &sync.Mutex{},
				queryBTL: func(ctx context.Context) (int64, error) {
					return 47, nil
				},
				queryBH: func(ctx context.Context) (int64, error) {
					return 10, nil
				},
			}
			btl, err := c.get(ctx)
			require.NoError(t, err, "must not return error")
			require.Equal(t, int64(47), c.lastAliveUntil, "must set last alive until")
			require.Equal(t, int64(47), btl, "must return alive until block height")
		})

		t.Run("with non-stale cache", func(t *testing.T) {
			c := &keepAliveCache{
				lastRefresh:         tm.Add(time.Second * -10),
				locker:              &sync.Mutex{},
				lastAliveUntil:      47,
				estimatedBlockSpeed: time.Second,
				lastBlockHeight:     10,
			}
			btl, err := c.get(ctx)
			require.NoError(t, err, "must not return error")
			require.Equal(t, int64(47), btl, "must return alive until block height")
		})
	})
}
