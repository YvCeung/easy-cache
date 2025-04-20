package easycache

import (
	"github.com/YvCeung/FastCache/lru"
	"sync"
)

// 封装了lru,在lru的基础上增加了并发属性

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	// 先上锁
	c.mu.Lock()

	//最后在释放锁
	defer c.mu.Unlock()
	if c.lru != nil {
		c.lru.Add(key, value)
	}
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
