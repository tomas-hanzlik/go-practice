# Cache

- Check `playground/play.go`

```go
config := types.CacheConfig{
	TTL:                      4,  // items time-to-live
	Capacity:                 100,// cache max size
	ExpCheckFrequency:        1,  // items ttl control
	GetAdaptersDataFrequency: 4,  // every 4 seconds
	AdaptersBufferSize: 	  0,  // unlimited
}
c := cache.NewCache(config)

// !BLOCKING! adapter to read from STDIN...
// - if PIPE -> read everything from it and then stop reading
// - if normal stdin -> take input from user and wait for command `STOP` to stop reading
// after reading ends -> continue with program execution
c.SetInputAdapter(cache.NewCommandLineInputAdapter(os.Stdin, 0))

// Generate 7 random items into the cache every 2 seconds
c.SetInputAdapter(cache.NewRandomInputAdapter(2, 7), 0)

c.AddItem(types.CacheItem{Key: "TEST1344", Value: "value"})
c.Dump("dumpster.txt")

```

- `func NewCache(config types.CacheConfig) *Cache`

- `func (cache *Cache) SetInputAdapter(adapter IAdapter)`

- `func (cache *Cache) CollectAdaptersData()`

- `func (cache *Cache) Size() int64`

- `func (cache *Cache) AddItem(item types.CacheItem)`

- `func (cache *Cache) GetItem(key string) (types.CacheItem, bool)`

- `func (cache *Cache) GetAllItems() *[]types.CacheItem`

- `func (cache *Cache) RemoveItem(key string)`

- `func (cache *Cache) RemoveAllItems()`

- `func (cache *Cache) RemoveExpiredItems()`

- `func (cache *Cache) Dump(filename string)`




## Cache config

```go
type CacheConfig struct {
	TTL                      int32 `json:"ttl"`                      // Expiration of items.
	Capacity                 int64 `json:"capacity"`                 // Capacity of the cache.
	ExpCheckFrequency        int32 `json:"expirationCheckFrequency"` // How often remove expired items. 0 to turn it off
	GetAdaptersDataFrequency int32 `json:"getAdaptersDataFrequency"` // How often we want to get data from adapters
	AdaptersBufferSize	     int64 `json:"adaptersBufferSize"`  // If we want to limit the amount of data before colleciton
}

```

## Adapters

- Cache collects data from adapters in specified intervals.

### RandomInputAdapter

 - Generates data in specified intervals 

### CommandLineAdapter

- Takes data from STDIN (/Pipe)
- Blocks all periodic tasks while taking data... to prevent messy stdout.


## Playground

- See `playgroung.play.go`

- You can run it like:
```
// Use Pipe to absorb data by CommandLineAdapter
`cat playground/test_data | go run playground/play.go`

// To type encoded keys in command line. Same format as in `playground/test_data`
`go run playground/play.go`
```

# REST API

## Endpoints

- Basic auth ... accounts in `.env`

- `GET     /ping`		  - ...
- `GET     /overview`     - cache state and configuration
- `GET     /cache`        - get all items
- `POST    /cache`		  - insert/upsert items
```
	> POST /cache HTTP/1.1
	> Content-Type: application/json

	| {
	| 	"data": [
	| 		{
	| 			"key": "TOMAS",
	| 			"value": "H"
	| 		}
	| 	]
	| }
```
- `DELETE  /cache`        - flush cache
- `GET     /cache/:key`   - get one item by key
- `DELETE  /cache/:key`   - delete one item by key


## Configuring API

- !!! Important to add `cert.pem` file for Cryptomood api... `cmd/app/cert.pem`  !!!
- You can customize settings in `cmd/app/.env`
```
DEBUG=1
ADAPTERS=random  				# `random` or `input` or `random,input`
CAPACITY=0 						# cache capacity
TTL=100							# cache items TTL
EXPIRATION_CHECK_FREQUENCY=10	# check and remove expired items from cache with frequency
GET_ADAPTERS_DATA_FREQUENCY=5	# collect items from adapters to cache with frequency
ADAPTERS_BUFFER_SIZE=10			# default size of buffers in adapters
ALLOWED_ACCOUNTS=1:1,2:2		# basic auth accounts
```

### Running in Docker

```sh

# Build image
docker build -t go-practice .

# Run container
docker run -p 8080:8080 -it go-practice

# Running on `0.0.0.0:8080`
```

### Makefile

- open it for more info 

- You can just use makefile to run tests/app/playground files

- Just used while personal development. If you want to use it via `docker` then add dependency on `make` in `Dockerfile`

