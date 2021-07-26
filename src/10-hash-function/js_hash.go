package hashutil

func JSHash(str []byte) uint32 {
	var hash uint32 = 1315423911

	for i := 0; i < len(str); i++ {
		hash ^= ((hash << 5) + uint32((rune(str[i]))) + (hash >> 2))
	}

	return hash
}