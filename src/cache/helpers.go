package cache

import (
	"sync"
	"time"
)

// Run periodic tasks and suspend if waitgroup blocked...
// !! Written to help suspend all periodic tasks while blocking input adapter is runnig.
func ExecutePeriodic(wg *sync.WaitGroup, frequency int32, f func()) {
	go func() {
		for {
			wg.Wait() // wait for some blocking actions
			select {
			case <-time.Tick(time.Duration(frequency) * time.Second):
				f()
			}
		}
	}()
}
