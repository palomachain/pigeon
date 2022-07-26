package collision

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	valsettypes "github.com/palomachain/pigeon/types/paloma/x/valset/types"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vizualni/whoops"
)

func TestCollisions(t *testing.T) {
	fakeErr := whoops.String("fake error")

	for _, tt := range []struct {
		name             string
		ctxdata          ctxdata
		setupPalomer     func(t *testing.T) *mockPalomer
		expErr           error
		allowedToExecute bool
	}{
		{
			name:    "when fetching snapshot returns an error it returns it back",
			ctxdata: ctxdata{},
			setupPalomer: func(t *testing.T) *mockPalomer {
				p := newMockPalomer(t)
				p.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(nil, fakeErr)
				return p
			},
			expErr: fakeErr,
		},
		{
			name:    "when fetching block height returns an error it returns it back",
			ctxdata: ctxdata{},
			setupPalomer: func(t *testing.T) *mockPalomer {
				p := newMockPalomer(t)
				p.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(nil, nil)
				p.On("BlockHeight", mock.Anything).Return(int64(0), fakeErr)
				return p
			},
			expErr: fakeErr,
		},
		{
			name:    "with a randomly crafted data it selects me for running the job",
			ctxdata: ctxdata{},
			setupPalomer: func(t *testing.T) *mockPalomer {
				p := newMockPalomer(t)
				p.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(&valsettypes.Snapshot{
					Validators: []valsettypes.Validator{
						{Address: sdk.ValAddress("me")},
						{Address: sdk.ValAddress("you")},
						{Address: sdk.ValAddress("they")},
					},
				}, nil)
				p.On("BlockHeight", mock.Anything).Return(int64(5), nil)
				return p
			},
			expErr:           nil,
			allowedToExecute: true,
		},
		{
			name:    "with a randomly crafted data it doesn selects me for running the job",
			ctxdata: ctxdata{},
			setupPalomer: func(t *testing.T) *mockPalomer {
				p := newMockPalomer(t)
				p.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(&valsettypes.Snapshot{
					Validators: []valsettypes.Validator{
						{Address: sdk.ValAddress("me")},
						{Address: sdk.ValAddress("you")},
						{Address: sdk.ValAddress("they")},
					},
				}, nil)
				p.On("BlockHeight", mock.Anything).Return(int64(10), nil)
				return p
			},
			expErr:           nil,
			allowedToExecute: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.setupPalomer(t)
			ctx, cleanup, err := GoStartLane(context.Background(), p, sdk.ValAddress("me"))
			require.ErrorIs(t, err, tt.expErr)
			if err != nil {
				return
			}

			require.NotNil(t, cleanup)
			t.Cleanup(cleanup)

			require.NotNil(t, ctx)

			require.Equal(t, tt.allowedToExecute, AllowedToExecute(ctx, []byte("random data")))
		})
	}
}

func TestPickingWinner(t *testing.T) {
	ctx := writeToContext(context.Background(), ctxdata{
		me: sdk.ValAddress("me"),
	})
}
