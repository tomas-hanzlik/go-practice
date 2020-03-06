# Cache

## Tests

- `go test -v ./...`

## Usage

- Check `playground/play.go`

## Config

```go
type CacheConfig struct {
	TTL               int32 // Expiration of items.
	Capacity          int64 // Capacity of the cache.
	ExpCheckFrequency int32 // How often remove expired items. 0 == never
	GetDataFrequency  int32 // How often we want to get data from adapters. 0 == never
}
```

## Adapters

- Cache collects data from adapters in specified intervals.

### RandomInputAdapter

 - Generates data in specified intervals 

### CommandLineAdapter

- Takes data from STDIN.
- Blocks all periodic tasks while taking data to prevent messy stdout.


## Playground

- See `playgroung.play.go`

- You can run it like:
```
// Use Pipe to absorb data by CommandLineAdapter
`cat playground/test_data | go run playground/play.go`

// To type encoded keys in command line. Same format as in `playground/test_data`
`go run playground/play.go`
```


## TODO

- setup with env variables
- simplify code
- add tests
- expose REST api (maybe go-swagger... will see)
- cryptomood
- makefile
