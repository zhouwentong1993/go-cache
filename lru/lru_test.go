package lru

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
	if ok, get := cache.Get("3"); !ok || string(get.(String)) != "1" {
		t.Fatalf("cache get error")
	}
	if ok, _ := cache.Get("222"); ok {
		t.Fatalf("cache get error")
	}

}
