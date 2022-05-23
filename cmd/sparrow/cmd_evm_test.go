package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvmKeys(t *testing.T) {
	t.TempDir()
	pass := "aaaa"

	rootCmd.SetIn(strings.NewReader(fmt.Sprintf("%s\n%s\n", pass, pass)))
	rootCmd.SetArgs([]string{
		"evm", "keys", "generate-new",
	})

	err := rootCmd.Execute()
	require.NoError(t, err)
}
