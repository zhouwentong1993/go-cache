package go_cache

import (
	"fmt"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	cache := New(2, nil)
	cache.Put("1", String("1"))
	fmt.Println(cache.currBytes)
	cache.Put("2", String("1"))
	fmt.Println(cache.currBytes)
	cache.Put("3", String("1"))
	if get, ok := cache.Get("3"); !ok || string(get.(String)) != "1" {
		t.Fatalf("cache get error")
	}
	if _, ok := cache.Get("222"); ok {
		t.Fatalf("cache get error")
	}

}
