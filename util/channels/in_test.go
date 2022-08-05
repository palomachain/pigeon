package channels

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFanIn(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	ch := FanIn(ch1, ch2, ch3)

	go func() {
		ch1 <- 1
		ch1 <- 1
		ch1 <- 1
		close(ch1)
	}()

	go func() {
		ch2 <- 1
		ch2 <- 1
		ch2 <- 1
		close(ch2)
	}()

	go func() {
		ch3 <- 1
		ch3 <- 1
		ch3 <- 1
		close(ch3)
	}()

	timeout := time.NewTicker(1 * time.Second)
	t.Cleanup(func() {
		timeout.Stop()
	})

	sum := 0
loop:
	for {
		select {
		case <-timeout.C:
			break loop
		case one, ok := <-ch:
			if !ok {
				break
			}
			sum += one
		}
	}

	require.Equal(t, 9, sum, "couldn't read all items from the channel")
}
