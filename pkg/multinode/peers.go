// 分布式环境下多节点选择
package multinode

// 类似于物流的中转站，分拨中心
type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

// 具体的快递员，负责送快递
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
