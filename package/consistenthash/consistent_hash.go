package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type ConsistentHash struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *ConsistentHash {

	m := &ConsistentHash{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

func (m *ConsistentHash) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			//构造虚拟节点
			hashCode := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hashCode)
			m.hashMap[hashCode] = key
		}
	}
	//哈希环进行排序，方便后续二分查找
	sort.Ints(m.keys)
}

func (m *ConsistentHash) GetNode(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hashCode := int(m.hash([]byte(key)))
	index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hashCode
	})

	//取个模就能满足环形的结构
	return m.hashMap[m.keys[index%len(m.keys)]]
}
