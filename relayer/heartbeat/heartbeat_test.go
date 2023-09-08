package heartbeat

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/VolumeFi/whoops"
	"github.com/palomachain/pigeon/relayer/mocks"
	"github.com/palomachain/pigeon/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHeartbeat(t *testing.T) {
	randErr := whoops.String("oh no")

	testdata := []struct {
		name        string
		setup       func(t *testing.T) (*mocks.PalomaClienter, context.Context)
		expectedErr error
	}{
		{
			name: "if querying alive returns an error, it does nothing",
			setup: func(t *testing.T) (*mocks.PalomaClienter, context.Context) {
				paloma := mocks.NewPalomaClienter(t)
				ctx := context.Background()
				ctx, cancel := context.WithCancel(ctx)
				cnt := 1
				paloma.On("QueryGetValidatorAliveUntilBlockHeight", mock.Anything).Return(int64(0), randErr).Times(3).Run(func(_ mock.Arguments) {
					cnt++
					if cnt == 3 {
						cancel()
					}
				})
				return paloma, ctx
			},
			expectedErr: randErr,
		},
		{
			name: "when validator is almost dead, it should call the keep alive method",
			setup: func(t *testing.T) (*mocks.PalomaClienter, context.Context) {
				paloma := mocks.NewPalomaClienter(t)
				ctx := context.Background()
				ctx, cancel := context.WithCancel(ctx)
				paloma.On("QueryGetValidatorAliveUntilBlockHeight", mock.Anything).Return(int64(100), nil).Times(1)
				paloma.On("BlockHeight", mock.Anything).Return(int64(90), nil).Times(1)
				paloma.On("KeepValidatorAlive", mock.Anything, "v1.4.0").Return(nil).Run(func(_ mock.Arguments) {
					cancel()
				})
				return paloma, ctx
			},
		},
		{
			name: "when validator is almost dead, but calling keep alive method returns an error, it should retry",
			setup: func(t *testing.T) (*mocks.PalomaClienter, context.Context) {
				paloma := mocks.NewPalomaClienter(t)
				ctx := context.Background()
				ctx, cancel := context.WithCancel(ctx)
				cnt := 1
				paloma.On("QueryGetValidatorAliveUntilBlockHeight", mock.Anything).Return(int64(100), nil).Times(1)
				paloma.On("BlockHeight", mock.Anything).Return(int64(90), nil).Times(1)
				paloma.On("KeepValidatorAlive", mock.Anything, "v1.4.0").Return(randErr).Times(3).Run(func(_ mock.Arguments) {
					cnt++
					if cnt == 3 {
						cancel()
					}
				})
				return paloma, ctx
			},
			expectedErr: randErr,
		},
		{
			name: "when validator still has time to live, it does not call keep alive method",
			setup: func(t *testing.T) (*mocks.PalomaClienter, context.Context) {
				paloma := mocks.NewPalomaClienter(t)
				ctx := context.Background()
				ctx, cancel := context.WithCancel(ctx)
				paloma.On("BlockHeight", mock.Anything).Return(int64(30), nil).Times(1)
				paloma.On("QueryGetValidatorAliveUntilBlockHeight", mock.Anything).Return(int64(100), nil).Times(1).Run(func(_ mock.Arguments) {
					cancel()
				})
				return paloma, ctx
			},
		},
	}

	asserter := assert.New(t)
	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			paloma := mocks.NewPalomaClienter(t)
			paloma, ctx := tt.setup(t)
			if ctx == nil {
				ctx = context.Background()
			}
			ctx, cancel := context.WithCancel(ctx)
			t.Cleanup(cancel)
			var locker testutil.FakeMutex
			h := New(
				paloma.QueryGetValidatorAliveUntilBlockHeight,
				paloma.BlockHeight,
				paloma.KeepValidatorAlive,
				30,
				"v1.4.0",
				&locker,
			)

			h.SetRetryFalloff(time.Millisecond)
			actualErr := h.Beat(ctx, &locker)

			asserter.Equal(tt.expectedErr, actualErr)
		})
	}
}

func TestLinearFalloffRetry(t *testing.T) {
	errTest := whoops.String("oh no")

	t.Run("when retrying fails, it returns an error", func(t *testing.T) {
		ctx := context.Background()
		locker := &sync.Mutex{}
		err := linearFalloffRetry(ctx, locker, "test", 3, time.Millisecond, func(ctx context.Context, locker sync.Locker) error {
			return errTest
		})
		require.Error(t, err, "must return error")
	})

	t.Run("when retrying succeeds, it does not return an error", func(t *testing.T) {
		ctx := context.Background()
		locker := &sync.Mutex{}
		cnt := 0
		times := make([]time.Time, 3)
		err := linearFalloffRetry(ctx, locker, "test", 3, time.Millisecond, func(ctx context.Context, locker sync.Locker) error {
			times[cnt] = time.Now()
			cnt++
			if cnt == 3 {
				return nil
			}
			return errTest
		})
		require.NoError(t, err, "must not return error")
		require.Equal(t, 3, cnt, "must have retried 3 times")
		for i := 0; i < 3; i++ {
			if i == 0 {
				continue
			}
			require.LessOrEqual(t, time.Duration(i)*time.Millisecond, times[i].Sub(times[i-1]), "must have waited at least i millisecond(s)")
		}
	})
}
