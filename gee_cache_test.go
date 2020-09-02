package go_cache

import (
	"testing"
)

func TestGetter(t *testing.T) {
	getterFunc := GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	get, _:= getterFunc.Get("key1")
	if string(get) != "key1" {
		t.Fatal("error")
	}
}
