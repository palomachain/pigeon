package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleUnrecoverableUsecases(t *testing.T) {
	baseErr := errors.New("hello world")

	unrecovorableErr := Unrecoverable(baseErr)

	require.True(t, IsUnrecoverable(unrecovorableErr))
	require.False(t, IsUnrecoverable(baseErr))
}
