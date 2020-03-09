package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	cache "tohan.net/go-practice/src/cache"
	types "tohan.net/go-practice/src/cache/types"
	crypto "tohan.net/go-practice/src/cryptomood"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const CryptomoodCertFile = "./cert.pem"
const CryptomoodServer = "internal-dev.api.cryptomood.com:30000"

const RandomInputAdapterInterval = 10
const RandomInputAdapterAmount = 10

type config struct {
	IsDebug                  bool     `env:"DEBUG"`
	Adapters                 []string `env:"ADAPTERS" envDefault:"" envSeparator:","`
	TTL                      int32    `env:"TTL" envDefault:"100"`
	Capacity                 int64    `env"CAPACITY" envDefault:"0"` // unlimited by default
	ExpirationCheckFrequency int32    `env:"EXPIRATION_CHECK_FREQUENCY" envDefault:"25"`
	GetAdaptersDataFrequency int32    `env:"GET_ADAPTERS_DATA_FREQUENCY" envDefault:"10"`
	AdaptersBufferSize		 int64	`env:"ADAPTERS_BUFFER_SIZE" envDefault:"0"` // unlimited by default
	AllowedAccounts			 []string `env:"ALLOWED_ACCOUNTS" envDefault:"" envSeparator:","`
}

func envConfig() *config {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &cfg
}

func initCache(cfg *config) *cache.Cache {
	config := types.CacheConfig{
		TTL:                      cfg.TTL,
		Capacity:                 cfg.Capacity,
		ExpCheckFrequency:        cfg.ExpirationCheckFrequency,
		GetAdaptersDataFrequency: cfg.GetAdaptersDataFrequency,
		AdaptersBufferSize:       cfg.AdaptersBufferSize,
	}
	c := cache.NewCache(config)

	// Set adapters.
	for _, adapterName := range cfg.Adapters {
		if adapterName == "input" {
			c.SetInputAdapter(cache.NewCommandLineInputAdapter(os.Stdin, cfg.AdaptersBufferSize))
		} else if adapterName == "random" {
			c.SetInputAdapter(cache.NewRandomInputAdapter(RandomInputAdapterInterval, RandomInputAdapterAmount, cfg.AdaptersBufferSize))
		}
	}
	return c
}

func initAPI(cfg *config, c *cache.Cache) *gin.Engine {
	// Configure API
	if cfg.IsDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	env := &CacheHandler{c}
	router := gin.Default()
	allowedAccounts :=gin.Accounts{}
	for _, accDetailsAsString := range cfg.AllowedAccounts {
		r := strings.Split(accDetailsAsString, ":") 
		allowedAccounts[r[0]] = r[1]
	}
	authorized := router.Group("/", gin.BasicAuth(allowedAccounts))
	{
		authorized.GET("/cache/", env.GetAllItems)
		authorized.POST("/cache/", env.AddItems)
		authorized.DELETE("/cache/", env.DeleteAllItems)
		authorized.GET("/cache/:key", env.GetItem)
		authorized.DELETE("/cache/:key", env.DeleteItem)
		authorized.GET("/overview", env.CacheOverview)
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "pong"})
	})

	return router
}



func main() {
	cfg := envConfig()
	c := initCache(cfg)

	// subscribe to sentiment API to and save records into the cache... 
	// not ideal implementation -> ASK
	go crypto.ConsumeSentiments(c, CryptomoodCertFile, CryptomoodServer)

	initAPI(cfg, c).Run(":8080")
}
