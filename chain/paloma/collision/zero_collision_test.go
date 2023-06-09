package collision

import (
	"context"
	"fmt"
	"testing"

	"github.com/VolumeFi/whoops"
	sdk "github.com/cosmos/cosmos-sdk/types"
	valsettypes "github.com/palomachain/paloma/x/valset/types"
	"github.com/palomachain/pigeon/util/slice"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
				p.On("QueryGetSnapshotByID", mock.Anything, uint64(0)).Return(&valsettypes.Snapshot{}, nil)
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

func TestEnsureThatSingleValidatorWillBePickedAtLeastOnce(t *testing.T) {
	validatorNum := 10
	randomJobsNum := 100

	validators := slice.IterN(validatorNum, func(index int) sdk.ValAddress {
		return sdk.ValAddress(fmt.Sprintf("%d", index))
	})

	// just ensure that the number of created validators was correct
	require.Equal(t, validatorNum, len(validators))

	winnerMapCount := make(map[string]int)
	for job := 0; job < randomJobsNum; job++ {
		winner := pickWinner([]byte("doesnt matter"), []byte(fmt.Sprintf("%d", job)), validators)
		winnerMapCount[winner.String()]++
	}

	// test that all validators were assigned a job at least once
	require.Equal(t, len(validators), len(winnerMapCount))
}
