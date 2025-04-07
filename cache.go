package FastCache

import (
	"github.com/YvCeung/FastCache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) Add(key string, value ByteView) {
	// 先上锁
	c.mu.Lock()

	//最后在释放锁
	defer c.mu.Unlock()
	if c.lru != nil {
		c.lru.Add(key, value)
	}
}
