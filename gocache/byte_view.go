package gocache

type ByteView struct {
	bytes []byte
}

func (b ByteView) Len() int {
	return len(b.bytes)
}

func (b ByteView) ToString() string {
	return string(b.bytes)
}

func (b ByteView) ByteSlice() []byte {
	return cloneBytes(b.bytes)
}

func cloneBytes(bs []byte) []byte {
	bytes := make([]byte, len(bs))
	copy(bytes, bs)
	return bytes
}
