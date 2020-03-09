package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	cache "tohan.net/go-practice/src/cache"
	types "tohan.net/go-practice/src/cache/types"
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
	data := gin.H{
		"config":              ch.cache.Config,
		"size":                ch.cache.Size(),
		"isUnlimitedCapacity": ch.cache.Config.Capacity == 0,
		"usedPercentage":      0,
	}
	if ch.cache.Config.Capacity > 0 {
		data["usedPercentage"] = int((100.0 / float64(ch.cache.Config.Capacity)) * float64(ch.cache.Size()))
	}
	ch.Resp(c, http.StatusOK, gin.H{"data": data})
}
