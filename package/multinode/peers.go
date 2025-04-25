// 分布式环境下多节点选择
package multinode

type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
