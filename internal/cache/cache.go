package cache

import (
	"fmt"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type Cache struct {
	cache *gocache.Cache
}

func NewCache() *Cache {
	return &Cache{
		cache: gocache.New(24*time.Hour, 1*time.Hour),
	}
}

func (c *Cache) IsValueTaken(key string, value string) bool {
	cacheKey := fmt.Sprintf("%s:%s:exists", key, value)

	_, found := c.cache.Get(cacheKey)

	return found
}

func (c *Cache) SetValue(key string, value string, exists bool) {
	cacheKey := fmt.Sprintf("%s:%s:exists", key, value)
	c.cache.Set(cacheKey, exists, 24*time.Hour)
}

func (c *Cache) DeleteField(key string, value string) {
	cacheKey := fmt.Sprintf("%s:%s:exists", key, value)
	c.cache.Delete(cacheKey)
}
