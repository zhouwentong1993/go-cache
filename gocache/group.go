package gocache

import "sync"

type Group struct {
	name      string
	mainCache cache
	getter    Getter
	picker    PeerPicker
}

var (
	mu     sync.Mutex
	groups = make(map[string]*Group)
)

func (g *Group) RegisterPeer(p PeerPicker) {
	g.picker = p
}

func (g *Group) Load(key string) (value ByteView, err error) {
	if peer, ok := g.picker.PickPeer(key); ok {
		data, peerErr := peer.Get(g.name, key)
		if peerErr != nil {
			return ByteView{bytes: nil}, peerErr
		}
		return ByteView{bytes: data}, nil
	} else {
		if getterValue, err := g.getter.Get(key); err == nil {
			g.mainCache.Add(key, ByteView{bytes: getterValue})
			return ByteView{bytes: cloneBytes(getterValue)}, nil
		} else {
			return ByteView{}, err
		}
	}
}

func New(name string, getter Getter, maxBytes uint64) *Group {
	if name == "" || getter == nil {
		panic("error input")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		mainCache: cache{cacheBytes: maxBytes},
		getter:    getter,
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	if name == "" {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	return groups[name]
}

func (g *Group) Get(key string) (bv ByteView, err error) {
	if key == "" {
		return ByteView{}, nil
	}
	if value, ok := g.mainCache.Get(key); ok {
		return value, nil
	} else {
		if value, err1 := g.getter.Get(key); err1 != nil {
			return ByteView{}, err1
		} else {
			view := ByteView{bytes: cloneBytes(value)}
			g.mainCache.Add(key, view)
			return view, nil
		}

	}
}
