package gocache

import (
	"fmt"
	"math"
	"testing"
)

var db = map[string]string{
	"zhangsan": "张三",
	"lisi":     "李四",
	"wangwu":   "王五",
}

func TestGetGroup(t *testing.T) {
	loadCount := make(map[string]int8, len(db))
	group := New("testGroup", GetterFunc(func(key string) ([]byte, error) {
		if v, ok := db[key]; ok {
			if _, ok := loadCount[key]; !ok {
				loadCount[key] = 0
			} else {
				loadCount[key] ++
			}
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}), math.MaxUint64)

	for k, v := range db {
		if get, err := group.Get(k); err != nil || string(get.bytes) != v {
			t.Fatal("load data error")
		}
		if _, err1 := group.Get(k); err1 != nil || loadCount[k] > 0 {
			t.Fatal("cache miss")
		}
	}

}
