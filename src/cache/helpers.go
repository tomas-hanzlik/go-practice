package cache

import (
	"sync"
	"time"
)

// Run periodic tasks and suspend if block in waitgroup...
func executePeriodic(wg *sync.WaitGroup, frequency int32, f func()) {
	go func() {
		for {
			wg.Wait() // wait for some blocking operations
			select {
			case <-time.Tick(time.Duration(frequency) * time.Second):
				f()
			}
		}
	}()
}
