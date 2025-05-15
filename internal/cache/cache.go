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

func (c *Cache) IsValueTaken(key string, value string) (bool, bool) {
	cacheKey := fmt.Sprintf("%s:%s:exists", key, value)

	if val, found := c.cache.Get(cacheKey); found {
		return val.(bool), true
	}

	return false, false
}

func (c *Cache) SetValue(key string, value string, exists bool) {
	cacheKey := fmt.Sprintf("%s:%s:exists", key, value)
	ttl := 24 * time.Hour
	if !exists {
		ttl = 10 * time.Minute
	}
	c.cache.Set(cacheKey, exists, ttl)
}

func (c *Cache) DeleteField(key string, value string) {
	cacheKey := fmt.Sprintf("%s:%s:exists", key, value)
	c.cache.Delete(cacheKey)
}
