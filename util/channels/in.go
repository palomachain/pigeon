package channels

import "sync"

func FanIn[T any](chans ...<-chan T) <-chan T {
	out := make(chan T, 1)
	var wg sync.WaitGroup
	wg.Add(len(chans))
	for _, ch := range chans {
		go func(ch <-chan T) {
			for data := range ch {
				out <- data
			}
			wg.Done()
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
