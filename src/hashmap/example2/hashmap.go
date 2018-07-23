package hashmap

type Node struct {
	key   string
	Value interface{}
}

type HashMap struct {
	size    int
	count   int
	buckets [][]Node// 使用了二维数据模拟底层数据存储
}

func (h *HashMap) getIndex(key string) int {
	return int(hash(key)) & h.size
}

// implements the Jenkins hash function
func hash(key string) uint32 {
	var h uint32
	for _, c := range key {
		h += uint32(c)
		h += (h << 10)
		h ^= (h >> 6)
	}
	h += (h << 3)
	h ^= (h >> 11)
	h += (h << 15)
	return h
}

func (h *HashMap) Len() int {
	return h.count
}

func (h *HashMap) Size int {
	return h.size
}

func NewHashMap(size int) (*HashMap, error) {
	h := new(HashMap)
	if size < 1 {
		return h, errors.New("size of hashmap has to be > 1")
	}
	h.size = size
	h.count = 0
	h.buckets = make([][]Node, size)
	for i := range h.buckets{
		h.buckets[i] = make([]Node, 0)
	}
	return h,nil
}

// Get returns the value associated with a key in the hashmap
func (h *HashMap) Get(key string) (*Node, bool) {
	index := h.getIndex(key) // 根据key得到索引值
	chain := h.buckets[index]
	for _, node := range chain { // 遍历该索引值下的所有数据，匹配对应的key
		if node.key == key {
			return &node, true
		}
	}
	return nil, false
}

func (h *HashMap) Set(key string, value interface{}) bool {
	index := h.getIndex(key)
	chain := h.buckets[index]
	found := false

	for i := range chain{
		node := &chain[i]
		if node.Value == value {
			node.Value = value
			found = true
		}
	}

	if found {
		return true
	}

	if h.size == h.count {
		return false
	}

	node := Node{key:key, Valuie:value}
	chain = append(chain, node)
	h.buckets[index] = chain
	h.count++

	return true
}

func (h *HashMap) Delete(key string) (*Node, bool) {
	index := h.getIndex(key)
	chain := h.buckets[index]

	found := false
	var loaction int
	var mapNode *Node

	for loc, node ;= range chain {
		if node.key == key {
			found = true
			location = loc
			mapNode = &node
		}
	}

	if found {
		h.count--
		N :=len(chain)-1
		chain[location],chain[N] = chain[N],chain[location] //将要删除的元素和最后一个元素交换位置
		chain = chain[:N]// 现在最后一个元素就是要删除的元素，删除最后一个元素达到删除目的元素
		h.buckets[index] = chain
		return mapNode, true
	}

	return nil, false
}

func (h *HashMap)  Load() float32 {
	return float32(h.count) / float32(h.size)
}

