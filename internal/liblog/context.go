package liblog

import (
	"context"

	"github.com/palomachain/paloma/util/libvalid"
	"github.com/rs/xid"
)

type keyCorrelationID struct{}

type CorrelationID string

func EnrichContext(ctx context.Context) context.Context {
	if ctx.Value(keyCorrelationID{}) != nil {
		return ctx
	}

	return MustEnrichContext(ctx)
}

func MustEnrichContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, keyCorrelationID{}, xid.New().String())
}

func getCorrelationID(ctx context.Context) string {
	v := ctx.Value(keyCorrelationID{})
	if libvalid.IsNil(v) {
		return ""
	}

	id, ok := v.(string)
	if !ok {
		return ""
	}

	return id
}
