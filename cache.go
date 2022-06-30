package cache

import (
	"time"
)

//The expiration should, and probably cannot in our case, be implemented asynchronously, 
// so think of a way to clean up the records when needed

type Cache struct {
	cache map[string]cacheValue
}

type cacheValue struct {
	value    string
	deadline time.Time
	shouldExpire bool
}

func NewCache() Cache {
	return Cache{
		cache: map[string]cacheValue{},
	}
}

/*`k, ok := Get(key string)` - returns the value associated with the key and the
boolean `ok` (true if exists, false if not),
if the deadline of the key/value pair has not been exceeded yet*/
func (c *Cache) Get(key string) (string, bool) {

	val, ok := c.cache[key]
	if !ok {
		return "", false
	}

	if val.shouldExpire {
		if time.Now().After(val.deadline) {
			// delete(c.cache, key)
			return "", false
		}
	}

	return val.value, true
}

/*`Put(key string, value string)` places a value with an associated key into cache.
Value put with this method never expired(have infinite deadline).
Putting into the existing key should overwrite the value*/
func (c Cache) Put(key, value string) {
	newCacheValue := cacheValue{
		value: value,
		shouldExpire: false,
	}
	c.cache[key] = newCacheValue

}

/*`Keys() []string` should return the slice of existing (non-expired keys)*/
func (c Cache) Keys() []string {
	keys := []string{}
	for key, val := range c.cache {
		if val.shouldExpire {
			if time.Now().After(val.deadline) {
				delete(c.cache, key)
			}else {
				keys = append(keys, key)
			}
		}else {
			keys = append(keys, key)
		}
	}
	return keys
}

// `PutTill(key string, value string, deadline time.Time)`
// Should do the same as `Put`, but the key/value pair should expire by given deadline
func (c Cache) PutTill(key, value string, deadline time.Time) {
	c.cache[key] = cacheValue{
		value:    value,
		deadline: deadline,
		shouldExpire: true,
	}
}
