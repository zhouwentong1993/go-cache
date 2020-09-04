package main

import (
	"fmt"
	"math"
	"net/http"
	."go-cache/gocache"
)

func main() {

	var db = map[string]string{
		"zhangsan": "张三",
		"lisi":     "李四",
		"wangwu":   "王五",
	}

	group := New("test", GetterFunc(func(key string) ([]byte, error) {
		if v, ok := db[key]; !ok {
			return nil, fmt.Errorf("no value")
		} else {
			return []byte(v), nil
		}
	}), math.MaxUint64)
	for k := range db {
		group.Get(k)
	}

	addr := "localhost:9999"
	handler := NewHTTPPool(addr)
	http.ListenAndServe(addr, handler)
}
