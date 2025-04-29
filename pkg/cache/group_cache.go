package cache

import (
	"fmt"
	"log"
	"sync"
)

type Group struct {
	name string
	//回调函数
	getter    Getter
	mainCache concurrentcache
}

type Getter interface {
	Get(key string) ([]byte, error)
}

// 自定义一个函数类型
type GetterFunc func(key string) ([]byte, error)

// 实现了Getter接口的函数 (接口型函数) 然后调用自己，常用作回调函数
func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("getter is nil")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:   name,
		getter: getter,
		mainCache: concurrentcache{
			cacheBytes: cacheBytes,
		},
	}

	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	group := groups[name]
	mu.RUnlock()
	return group
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[FastCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: cloneBytes(bytes)}

	//放入缓存
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
