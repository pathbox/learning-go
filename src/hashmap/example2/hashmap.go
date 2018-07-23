package hashmap

type Node struct {
	key   string
	Value interface{}
}

type HashMap struct {
	size    int
	count   int
	buckets [][]Node
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
