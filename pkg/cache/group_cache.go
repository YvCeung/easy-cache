package cache

import (
	"fmt"
	"github.com/YvCeung/easy-cache/pkg/multinode"
	"log"
	"sync"
)

type Group struct {
	name string
	//回调函数
	getter    Getter
	mainCache concurrentcache

	//集成多节点的能力
	peerPicker multinode.PeerPicker
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
	if g.peerPicker != nil {
		if peer, ok := g.peerPicker.PickPeer(key); ok {
			data, err := g.getFromPeer(peer, key)
			if err == nil {
				return data, nil
			}
			//从别的节点没有获取到数据
			log.Println("[easycache] Failed to get from peer ", peer)
		}
	}
	//走本地
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

func (g *Group) RegisterPeerPicker(peerPicker multinode.PeerPicker) {
	if g.peerPicker != nil {
		panic("RegisterPeerPicker calles more than once")
	}
	g.peerPicker = peerPicker
}

func (g *Group) getFromPeer(peerGetter multinode.PeerGetter, key string) (ByteView, error) {
	data, err := peerGetter.Get(g.name, key)

	if err != nil {
		return ByteView{}, nil
	}
	return ByteView{b: data}, nil
}
