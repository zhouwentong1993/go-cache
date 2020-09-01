package go_cache

type Getter interface {
	Get(key string) (error, []byte)
}

type GetterFunc func(key string) (error, []byte)

func (f GetterFunc) Get(key string) (error, []byte) {
	return f(key)
}
