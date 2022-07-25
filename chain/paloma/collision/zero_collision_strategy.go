package collision

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/palomachain/pigeon/util/slice"
	log "github.com/sirupsen/logrus"
)

// Zero collision strategy ensures that all pigeons are using the same way of determinig if they can execute a message or not.
// That way we can distribute the jobs without having to write data to paloma so that paloma could do it, as it could be slow.
// Of course, paloma does not care about the strategy pigeons are using. Somebody could rewrite pigeon to do send all the messages,
// but then they would be "fighting" with paloma and they (and other pigeons) would spend gas unceccecary.
// A "bad" actor would only be doing a job for us :).

const (
	round         = 10
	tickerTimeout = 5 * time.Second
)

type ctxKeyType int

var ctxKey ctxKeyType = 1

type ctxdata struct {
	me          sdk.ValAddress
	valset      []sdk.ValAddress // ordered!
	blockHeight int64
}

type palomer interface {
	BlockHeight(context.Context) (int64, error)
	GetLatestValset(context.Context)
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
		data.valset,
	)

	if winner.Equals(data.me) {
		return true
	}

	return false
}

func GoStartLane(ctx context.Context, p palomaer, me sdk.ValAddress) (context.Context, func(), error) {
	valset, err := p.GetLatestValset(ctx)
	if err != nil {
		return nil, nil, err
	}

	blockHeight, err := p.BlockHeight(ctx)
	if err != nil {
		return nil, nil, err
	}

	// rounds down the block height to the nearest 10
	// 52 -> 50
	// 119 -> 110
	blockHeight = roundBlockHeight(blockHeight)

	ctx = context.WithValue(ctx, ctxKey, ctxdata{
		me:          me,
		valset:      valset,
		blockHeight: blockHeight,
	})

	ctx, cancelCtx := context.WithCancel(ctx)

	cancel := func() {
		cancelCtx()
	}

	go func() {
		defer cancelCtx()

		ticker := time.NewTicker(tickerTimeout)
		defer ticker.Stop()

		for range ticker.C {
			if ctx.Err() != nil {
				return
			}
			newBlockHeight, err := p.BlockHeight(ctx)
			if err != nil {
				log.WithError(err).Error("error getting block height")
				return
			}
			newBlockHeight = roundBlockHeight(newBlockHeight)
			if newBlockHeight != blockHeight {
				return
			}
		}
	}()

	return ctx, cancel, nil
}

func pickWinner(seed []byte, dump []byte, vals []sdk.ValAddress) sdk.ValAddress {
	h := fnv.New64()
	h.Write(append(seed[:], dump...))
	hash := h.Sum64() % uint64(len(vals))

	vals = slice.Filter(vals, func(val sdk.ValAddress) bool {
		h := fnv.New64()
		h.Write(val)
		valhash := h.Sum64() % uint64(len(vals))
		return valhash == hash
	})

	if len(vals) <= 0 {
		return nil
	}

	if len(vals) > 1 {
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(h.Sum64()))
		dump = append(dump, b...)
		return pickWinner(seed, dump, vals)
	}

	return vals[0]
}

func roundBlockHeight(bh int64) int64 {
	return bh - (bh % round)
}
