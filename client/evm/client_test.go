package evm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadingStoredContracts(t *testing.T) {
	t.Run("it successfully reads the hello.json contract", func(t *testing.T) {
		c := StoredContracts()
		require.GreaterOrEqual(t, len(c), 1)
		require.Contains(t, c, "hello")
		fmt.Println(c["hello"])
	})
}
