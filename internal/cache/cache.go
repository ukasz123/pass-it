package cache

import (
	"time"

	ttlcache "github.com/jellydator/ttlcache/v3"
)

const DefaultTTL = 5 * time.Minute

type Item[Value any] struct {
	val Value
}

type Cache[Value any] interface {
	/// Sets a value in the cache.
	Set(key string, value Value)
	/// Returns nil if not found.
	Get(key string) *Item[Value]
	/// Returns true when the key is in the cache.
	Contains(key string) bool
}

func NewCache[Value any]() Cache[Value] {
	return &cache[Value]{
		ttlcache: ttlcache.New(
			ttlcache.WithTTL[string, Value](DefaultTTL),
			ttlcache.WithDisableTouchOnHit[string, Value](),
		),
	}
}

type cache[Value any] struct {
	ttlcache *ttlcache.Cache[string, Value]
}

func (c *cache[Value]) Set(key string, value Value) {
	c.ttlcache.Set(key, value, ttlcache.DefaultTTL)
}

func (c *cache[Value]) Get(key string) *Item[Value] {
	item, found := c.ttlcache.GetAndDelete(key)
	if !found {
		return nil
	}
	v := item.Value()
	return &Item[Value]{val: v}
}

func (c *cache[Value]) Contains(key string) bool {
	if v:=c.ttlcache.Get(key); v != nil {
		return !v.IsExpired()
	}
	return false
}
func (i *Item[Value]) Value() Value {
	return i.val
}
