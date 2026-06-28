//package memory
//
//import (
//	"sync"
//	"time"
//)
//
//type UserCache struct {
//	Email      string `json:"email"`
//	Password   string `json:"password"`
//	OTP        string `json:"otp"`
//	MyDuration time.Time
//	CreatedAt  time.Time
//}
//type MemoryCache struct {
//	data map[string]UserCache
//	mu   sync.Mutex
//}
//
//func NewCache(data map[string]UserCache, mu sync.Mutex) MemoryCache {
//	return MemoryCache{
//		data: data,
//		mu:   mu,
//	}
//}
//
//func (m *MemoryCache) Set(email string, user UserCache) {
//	m.mu.Lock()
//	defer m.mu.Unlock()
//	m.data[email] = user
//}
//
//func (m *MemoryCache) Get(email string) (UserCache, bool) {
//	m.mu.Lock()
//	defer m.mu.Unlock()
//	user, ok := m.data[email]
//	if !ok {
//		return UserCache{}, false
//	}
//	return user, true
//}
//
//func (m *MemoryCache) Delete(email string) {
//	m.mu.Lock()
//	defer m.mu.Unlock()
//	delete(m.data, email)
//}

package memory

import (
	"time"

	memoryCache "github.com/patrickmn/go-cache"
)

type cache struct {
	memory *memoryCache.Cache
}

type MemoryCache interface {
	Set(key string, value any, d time.Duration)
	Get(key string) (any, bool)
	Replace(key string, value any, ttl time.Duration) error
	Remove(key string)
}

func NewMemoryCache() MemoryCache {
	c := memoryCache.New(memoryCache.NoExpiration, time.Minute)

	return &cache{memory: c}
}

func (c *cache) Set(key string, value any, ttl time.Duration) {
	c.memory.Set(key, value, ttl)
}

func (c *cache) Get(key string) (any, bool) {
	return c.memory.Get(key)
}

func (c *cache) Replace(key string, value any, ttl time.Duration) error {
	return c.memory.Replace(key, value, ttl)
}

func (c *cache) Remove(key string) {
	c.memory.Delete(key)
}
