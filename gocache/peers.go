package gocache

import "go-cache/gocache/pb"

// 分布式节点的抽象
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 节点获取数据的抽象
type PeerGetter interface {
	Get(req *pb.Request, resp *pb.Response) (err error)
}
