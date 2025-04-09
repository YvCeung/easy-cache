package lru

import "testing"

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
