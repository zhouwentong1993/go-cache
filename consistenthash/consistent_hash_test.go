package consistenthash

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMap_AddNodes(t *testing.T) {
	keys := []int{1, 2, 3, 4, 5, 6,}
	for i, key := range keys {
		if key == 2 || key == 4 {
			keys = append(keys[:i], keys[i:]...)
		}
	}
	fmt.Println(keys)

	m := New(3, func(data []byte) uint32 {
		atoi, _ := strconv.Atoi(string(data))
		return uint32(atoi)
	})
	m.AddNodes("3", "6", "9")
	fmt.Println("Map 现有节点：", m.keys)
	// 40 的下一个节点是 60
	if m.GetNode("40") != "6" {
		t.Fatal("节点错误")
	}
	m.DeleteNode("6")
	// 删除之后应该落到 9
	if m.GetNode("40") != "9" {
		t.Fatal("节点错误")
	}
	m.AddNodes("8")
	// 加了 8 之后，应该到 8 上
	if m.GetNode("40") != "8" {
		t.Fatal("节点错误")
	}
}
