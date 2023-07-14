package relayer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRelayer_SetAppVersion(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "sets app version correctly",
			input:    "v1.4.0",
			expected: "v1.4.0",
		},
	}

	asserter := assert.New(t)

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			r := Relayer{}
			asserter.Equal("", r.appVersion)

			r.SetAppVersion(tt.input)
			asserter.Equal(tt.expected, r.appVersion)
		})
	}
}
