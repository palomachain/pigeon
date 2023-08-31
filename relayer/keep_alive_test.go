package relayer

import (
	"context"
	"testing"
	"time"

	"github.com/VolumeFi/whoops"
	"github.com/palomachain/pigeon/relayer/mocks"
	"github.com/palomachain/pigeon/testutil"
	timemock "github.com/palomachain/pigeon/util/time/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestKeepAlive(t *testing.T) {
	randErr := whoops.String("oh no")

	testdata := []struct {
		name        string
		setup       func(t *testing.T) (*mocks.PalomaClienter, *timemock.Time, context.Context)
		expectedErr error
	}{
		{
			name: "if querying alive returns an error, it does nothing",
			setup: func(t *testing.T) (*mocks.PalomaClienter, *timemock.Time, context.Context) {
				paloma := mocks.NewPalomaClienter(t)
				tm := timemock.NewTime(t)
				ctx := context.Background()
				ctx, cancel := context.WithCancel(ctx)
				paloma.On("QueryGetValidatorAliveUntilBlockHeight", mock.Anything).Return(int64(0), randErr).Times(1).Run(func(_ mock.Arguments) {
					cancel()
				})
				return paloma, tm, ctx
			},
			expectedErr: randErr,
		},
		{
			name: "when validator is almost dead, it should call the keep alive method",
			setup: func(t *testing.T) (*mocks.PalomaClienter, *timemock.Time, context.Context) {
				paloma := mocks.NewPalomaClienter(t)
				tm := timemock.NewTime(t)
				ctx := context.Background()
				ctx, cancel := context.WithCancel(ctx)
				paloma.On("QueryGetValidatorAliveUntilBlockHeight", mock.Anything).Return(int64(100), nil).Times(1)
				paloma.On("BlockHeight", mock.Anything).Return(int64(90), nil).Times(1)
				paloma.On("KeepValidatorAlive", mock.Anything, "v1.4.0").Return(nil).Run(func(_ mock.Arguments) {
					cancel()
				})
				return paloma, tm, ctx
			},
		},
		{
			name: "when validator is almost dead, but calling keep alive method returns an error it does nothing",
			setup: func(t *testing.T) (*mocks.PalomaClienter, *timemock.Time, context.Context) {
				paloma := mocks.NewPalomaClienter(t)
				tm := timemock.NewTime(t)
				ctx := context.Background()
				ctx, cancel := context.WithCancel(ctx)
				paloma.On("QueryGetValidatorAliveUntilBlockHeight", mock.Anything).Return(int64(100), nil).Times(1)
				paloma.On("BlockHeight", mock.Anything).Return(int64(90), nil).Times(1)
				paloma.On("KeepValidatorAlive", mock.Anything, "v1.4.0").Return(randErr).Run(func(_ mock.Arguments) {
					cancel()
				})
				return paloma, tm, ctx
			},
			expectedErr: randErr,
		},
		{
			name: "when validator still has time to live, it does not call keep alive method",
			setup: func(t *testing.T) (*mocks.PalomaClienter, *timemock.Time, context.Context) {
				paloma := mocks.NewPalomaClienter(t)
				tm := timemock.NewTime(t)
				ctx := context.Background()
				ctx, cancel := context.WithCancel(ctx)
				paloma.On("BlockHeight", mock.Anything).Return(int64(30), nil).Times(1)
				paloma.On("QueryGetValidatorAliveUntilBlockHeight", mock.Anything).Return(int64(100), nil).Times(1).Run(func(_ mock.Arguments) {
					cancel()
				})
				return paloma, tm, ctx
			},
		},
	}

	asserter := assert.New(t)
	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			tm := timemock.NewTime(t)
			paloma := mocks.NewPalomaClienter(t)
			paloma, tm, ctx := tt.setup(t)
			if ctx == nil {
				ctx = context.Background()
			}
			ctx, cancel := context.WithCancel(ctx)
			t.Cleanup(cancel)
			r := Relayer{
				relayerConfig: Config{
					KeepAliveLoopTimeout:    50 * time.Millisecond,
					KeepAliveBlockThreshold: 30,
				},
				time:         tm,
				palomaClient: paloma,
				appVersion:   "v1.4.0",
			}

			var locker testutil.FakeMutex
			actualErr := r.keepAlive(ctx, &locker)

			asserter.Equal(tt.expectedErr, actualErr)
		})
	}
}
