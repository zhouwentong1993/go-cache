package gocache

// 分布式节点的抽象
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 节点获取数据的抽象
type PeerGetter interface {
	Get(group string, key string) (data []byte, err error)
}
