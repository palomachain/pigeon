package rotator_test

import (
	"context"
	"sync"
	"testing"

	"github.com/palomachain/pigeon/util/rotator"
	"github.com/stretchr/testify/assert"
)

func Test_Rotator_RotateKeys(t *testing.T) {
	exp := []string{"foo", "bar", "baz"}
	expected := ""
	idx := 0
	fn := func(s string) {
		idx++
		expected = exp[idx%3]
		assert.Equal(t, expected, s)
	}

	r := rotator.New(fn, exp...)
	wg := &sync.WaitGroup{}

	// Confirm rotator is concurrency safe
	wg.Add(3 * 99)
	for i := 0; i < 3; i++ {
		go func() {
			for i := 1; i < 100; i++ {
				wg.Add(1)
				r.RotateKeys(context.Background())
				wg.Done()
			}
		}()
	}
}
