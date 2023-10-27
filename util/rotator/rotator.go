package rotator

import (
	"container/ring"
	"context"
	"sync"
)

type (
	// Assigner is a function that's called after the keys are rotated. Use it to set external dependencies.
	Assigner func(s string)

	// Rotator is a struct that holds the keys and the function to call after the keys are rotated.
	Rotator struct {
		fn   Assigner
		ring *ring.Ring
		m    sync.Mutex
	}
)

// New creates a new Rotator instance. The rotateFn is called with the next
// key name when RotateKeys is called.
func New(rotateFn Assigner, keyNames ...string) *Rotator {
	ring := ring.New(len(keyNames))
	for _, v := range keyNames {
		ring.Value = v
		ring = ring.Next()
	}

	return &Rotator{
		fn:   rotateFn,
		ring: ring,
		m:    sync.Mutex{},
	}
}

// RotateKeys rotates the keys and calls the rotateFn with the next key name.
func (r *Rotator) RotateKeys(ctx context.Context) {
	r.m.Lock()
	defer r.m.Unlock()

	r.ring = r.ring.Next()
	r.fn(r.ring.Value.(string))
}
