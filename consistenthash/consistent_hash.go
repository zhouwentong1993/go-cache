package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           // 使用的 Hash 算法
	replicas int            // key 散列的个数
	keys     []int          // 所有的 key = 原始 key 个数 * key 散列的个数
	hashmap  map[int]string // 保存虚拟节点与原始节点关系
}

func New(replicas int, hash Hash) *Map {
	if hash == nil {
		hash = crc32.ChecksumIEEE
	}
	return &Map{
		hash:     hash,
		replicas: replicas,
		hashmap:  make(map[int]string),
	}
}

// 新增节点
func (m *Map) AddNodes(nodes ...string) {
	for _, key := range nodes {
		// 创建虚拟节点
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(key + strconv.Itoa(i))))
			// 保留映射关系
			m.keys = append(m.keys, hash)
			m.hashmap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) GetNode(nodeName string) string {
	hash := int(m.hash([]byte(nodeName)))
	index := sort.Search(len(m.keys), func(i int) bool {
		return hash <= m.keys[i]
	})

	return m.hashmap[m.keys[index%len(m.keys)]]
}

func (m *Map) DeleteNode(nodeNames ...string) {
	for _, name := range nodeNames {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(name + strconv.Itoa(i))))
			// 删除在 keys 中的数据
			m.deleteKey(hash)
			// 删除在 map 中的数据
			delete(m.hashmap, hash)
		}
	}
}

func (m *Map) deleteKey(hash int) {
	for idx, k := range m.keys {
		if hash == k {
			m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
			return
		}
	}
}
