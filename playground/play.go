package main

import (
	"os"
	"time"

	cache "tohan.net/go-practice/src/cache"
	types "tohan.net/go-practice/src/cache/types"
)

func main() {
	// USAGE:
	// - `cat playground/test_data | go run playground/play.go`
	// - `go run playground/play.go`

	config := types.CacheConfig{
		TTL:                      4,
		Capacity:                 100,
		ExpCheckFrequency:        1,
		GetAdaptersDataFrequency: 4,
	}
	c := cache.NewCache(config)

	// - if PIPE -> read everything from it and then stop reading
	// - if normal stdin -> take input from user and wait for command `STOP` to stop reading
	c.SetInputAdapter(cache.NewCommandLineInputAdapter(os.Stdin, 0))

	// Generate 7 random items into the cache every 2 seconds
	c.SetInputAdapter(cache.NewRandomInputAdapter(2, 7, 0))

	time.Sleep(10 * time.Second)

	c.AddItem(types.CacheItem{Key: "TEST1344", Value: "value"})
	c.Dump("dumpster.txt")
}
