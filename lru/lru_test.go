package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("value1"))
	if v, ok := lru.Get("key1"); !ok || v.(String) != "value1" {
		// Fatalf方法作用：1、输出格式化的错误信息 2、立即终止当前测试函数
		t.Fatalf("cache hit key1=value1 failed")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestAdd(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key", String("1"))
	lru.Add("key", String("111"))

	if lru.nBytes != int64(len("key")+len("111")) {
		t.Fatal("expected 6 but got", lru.nBytes)
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	// 放置第三个元素的时候就会把第一个元素给删除了
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

/*
测试 LRU 缓存中，当触发淘汰策略时，是否会正确地调用传入的回调函数 OnEvicted 并记录被淘汰的键。
*/
func TestOnEvicted(t *testing.T) {
	//定义了一个空的 keys 切片，用于记录被淘汰的 key
	keys := make([]string, 0)

	//创建了一个回调函数 callback，每当有 key 被淘汰，就追加进 keys 数组
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}

	//创建一个最大容量为 10 字节的 LRU 缓存 ,把刚刚定义的 callback 注册进去，供淘汰时调用
	lru := New(int64(10), callback)

	// 10 个字节
	lru.Add("key1", String("123456"))

	//添加完之后大于10，则会触发删除
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))

	//添加完之后大于10，则会触发删除
	lru.Add("k4", String("k4"))

	//按照最久未使用原则，会先淘汰 "key1"，再淘汰 "k2"
	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
