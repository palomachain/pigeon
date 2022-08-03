// Zero collision strategy ensures that all pigeons are using the same way of determinig if they can execute a message or not.
// That way we can distribute the jobs without having to write data to paloma so that paloma could do it, as it could be slow.
// Of course, paloma does not care about the strategy pigeons are using. Somebody could rewrite pigeon to do send all the messages,
// but then they would be "fighting" with paloma and they (and other pigeons) would spend gas unceccecary.
// A "bad" actor would only be doing a job for us :).
package collision

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/palomachain/pigeon/chain"
	valset "github.com/palomachain/pigeon/types/paloma/x/valset/types"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

const (
	round         = 10
	tickerTimeout = 5 * time.Second
)

type ctxKeyType int

var ctxKey ctxKeyType = 1

type ctxdata struct {
	me          sdk.ValAddress
	validators  []sdk.ValAddress // ordered!
	blockHeight int64
}

//go:generate mockery --name=palomer --inpackage --testonly
type palomer interface {
	BlockHeight(context.Context) (int64, error)
	QueryGetSnapshotByID(ctx context.Context, id uint64) (*valset.Snapshot, error)
}

func AllowedToExecute(ctx context.Context, dump []byte) bool {
	rawdata := ctx.Value(ctxKey)
	if rawdata == nil {
		panic("data about collision detection is not stored in the context")
	}

	data := rawdata.(ctxdata)

	winner := pickWinner(
		[]byte(fmt.Sprintf("%d", data.blockHeight)),
		dump,
		data.validators,
	)

	if winner.Equals(data.me) {
		return true
	}

	return false
}

func writeToContext(ctx context.Context, data ctxdata) context.Context {
	return context.WithValue(ctx, ctxKey, data)
}

func GoStartLane(ctx context.Context, p palomer, me sdk.ValAddress) (context.Context, func(), error) {
	snapshot, err := p.QueryGetSnapshotByID(ctx, 0)
	if err != nil {
		return nil, nil, err
	}
	if snapshot == nil {
		return nil, nil, chain.ErrNotFound.WrapS("snapshot does not exist")
	}
	if err != nil {
		return nil, nil, err
	}

	blockHeight, err := p.BlockHeight(ctx)
	if err != nil {
		return nil, nil, err
	}

	blockHeight = roundBlockHeight(blockHeight)

	ctx = writeToContext(ctx, ctxdata{
		me:          me,
		validators:  slice.Map(snapshot.Validators, func(v valset.Validator) sdk.ValAddress { return v.Address }),
		blockHeight: blockHeight,
	})

	ctx, cancelCtx := context.WithCancel(ctx)

	go func() {
		defer cancelCtx()

		ticker := time.NewTicker(tickerTimeout)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				newBlockHeight, err := p.BlockHeight(ctx)
				if err != nil {
					log.WithError(err).Error("error getting block height")
					return
				}
				newBlockHeight = roundBlockHeight(newBlockHeight)
				if newBlockHeight != blockHeight {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// calling the cancelCtx will close the goroutine from above as it's
	// checking on the context.Done channel
	return ctx, cancelCtx, nil
}

func pickWinner(seed []byte, dump []byte, vals []sdk.ValAddress) sdk.ValAddress {
	if len(vals) == 0 {
		return nil
	}
	datahash := fnv.New64()
	datahash.Write(append(seed[:], dump...))

	for {
		possibleWinners := slice.Filter(vals, func(val sdk.ValAddress) bool {
			h := fnv.New64()
			h.Write(val)
			valhash := (datahash.Sum64() + h.Sum64()) % uint64(len(vals))
			return valhash == 0
		})
		if len(possibleWinners) != 1 {
			b := make([]byte, 8)
			binary.BigEndian.PutUint64(b, uint64(datahash.Sum64()))
			datahash.Write(b)
			continue
		}
		return possibleWinners[0]
	}

}

// rounds down the block height to the nearest 10
// 52 -> 50
// 119 -> 110
func roundBlockHeight(bh int64) int64 {
	return bh - (bh % round)
}
