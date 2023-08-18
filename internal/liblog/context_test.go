package liblog

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func TestEnrichContext(t *testing.T) {
	t.Run("with prepopulated context", func(t *testing.T) {
		ctx := EnrichContext(context.Background())
		id := getCorrelationID(ctx)
		r := EnrichContext(ctx)
		require.Equal(t, id, getCorrelationID(r), "should not override existing correlation ID")
	})

	t.Run("with empty context", func(t *testing.T) {
		r := EnrichContext(context.Background())
		id := getCorrelationID(r)
		require.NotEmpty(t, id)
		require.Len(t, id, 20)
		_, err := xid.FromString(id)
		require.NoError(t, err)
	})
}

func TestMustEnrichContext(t *testing.T) {
	t.Run("with prepopulated context", func(t *testing.T) {
		ctx := EnrichContext(context.Background())
		id := getCorrelationID(ctx)
		r := MustEnrichContext(ctx)
		require.NotEqual(t, id, getCorrelationID(r), "should override existing correlation ID")
	})

	t.Run("with empty context", func(t *testing.T) {
		r := MustEnrichContext(context.Background())
		id := getCorrelationID(r)
		require.NotEmpty(t, id)
		require.Len(t, id, 20)
		_, err := xid.FromString(id)
		require.NoError(t, err)
	})
}

func TestGetCorrelationID(t *testing.T) {
	t.Run("with empty context", func(t *testing.T) {
		id := getCorrelationID(context.Background())
		require.Empty(t, id)
	})

	t.Run("with id type mismatch", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), keyCorrelationID{}, 42)
		id := getCorrelationID(ctx)
		require.Empty(t, id)
	})

	t.Run("with populated context", func(t *testing.T) {
		ctx := EnrichContext(context.Background())
		id := getCorrelationID(ctx)
		require.NotEmpty(t, id)
		require.Len(t, id, 20)
		_, err := xid.FromString(id)
		require.NoError(t, err)
	})
}
