package cache

import (
	"time"
)

// Execute function periodically
func executePeriodic(frequency int32, f func()) {
	func() {
		for {
			select {
			case <-time.Tick(time.Duration(frequency) * time.Second):
				f()
			}
		}
	}()
}
