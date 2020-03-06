package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	cache "../../src/cache"
	types "../../src/cache/types"
)

type BulkInsert struct {
	Data []types.CacheItem `json:"data"`
}
type CacheHandler struct {
	cache *cache.Cache
}

func (ch *CacheHandler) Resp(c *gin.Context, status int, resp gin.H) {
	resp["status"] = status
	c.JSON(status, resp)
}

func (ch *CacheHandler) GetAllItems(c *gin.Context) {
	items := ch.cache.GetAllItems()
	ch.Resp(c, http.StatusOK, gin.H{"data": items, "count": len(*items)})
}
func (ch *CacheHandler) DeleteAllItems(c *gin.Context) {
	ch.cache.RemoveAllItems()
	ch.Resp(c, http.StatusOK, gin.H{})
}

func (ch *CacheHandler) AddItems(c *gin.Context) {
	var bulkInsert BulkInsert
	if err := c.ShouldBindJSON(&bulkInsert); err != nil {
		ch.Resp(c, http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	for _, item := range bulkInsert.Data {
		ch.cache.AddItem(item)
	}
	ch.Resp(c, http.StatusCreated, gin.H{"message": "Added successfuly", "count": len(bulkInsert.Data)})
}

func (ch *CacheHandler) GetItem(c *gin.Context) {
	item, ok := ch.cache.GetItem(c.Param("key"))
	if !ok {
		ch.Resp(c, http.StatusNotFound, gin.H{"message": "Item not found"})
		return
	}
	ch.Resp(c, http.StatusOK, gin.H{"data": item})
}

func (ch *CacheHandler) DeleteItem(c *gin.Context) {
	ch.cache.RemoveItem(c.Param("key"))
	ch.Resp(c, http.StatusOK, gin.H{})
}

func (ch *CacheHandler) CacheOverview(c *gin.Context) {
	ch.Resp(c, http.StatusOK, gin.H{
		"data": gin.H{"config": ch.cache.Config, "size": ch.cache.Size()},
	})
}

func main() {
	config := types.CacheConfig{
		TTL:               10,
		Capacity:          10000,
		ExpCheckFrequency: 1,
		GetDataFrequency:  1,
	}
	c := cache.NewCache(config)

	// Generate 150 random items into the cache every 5 seconds
	c.SetInputAdapter(cache.NewRandomInputAdapter(5, 150))
	env := &CacheHandler{c}

	router := gin.Default()
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"1":      "1",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))
	{
		authorized.GET("/cache/", env.GetAllItems)
		authorized.POST("/cache/", env.AddItems)
		authorized.DELETE("/cache/", env.DeleteAllItems)
		authorized.GET("/cache/:key", env.GetItem)
		authorized.DELETE("/cache/:key", env.DeleteItem)
		authorized.GET("/overview", env.CacheOverview)
	}

	router.GET("/swagger", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello")
	})

	router.Run(":8080")
}
