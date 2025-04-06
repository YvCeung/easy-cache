package lru

import "container/list"

type Value interface {
	Len() int
}

type Cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// 构造函数
func (cache *Cache) New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		nBytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 增加或修改元素
func (c *Cache) Add(key string, v Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(v.Len()) - int64(kv.value.Len())
		kv.value = v
	} else {
		element := c.ll.PushFront(&entry{key: key, value: v})
		c.cache[key] = element
		c.nBytes += int64(v.Len())
	}

	for c.maxBytes > 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// 获取元素
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		//移动到队尾
		c.ll.MoveToFront(element)
		//获取真正的Entry节点
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	back := c.ll.Back()
	if back != nil {
		//删除元素
		c.ll.Remove(back)
		kv := back.Value.(*entry)
		//清空字典
		delete(c.cache, kv.key)

		//更新容量
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())

		//执行回调
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}
