package hashutil

func BPHash(str []byte) uint32 {
	var hash uint32 = 0
	for i := 0; i < len(str); i++ {
		hash = (hash << 7) ^ uint32(rune(str[i]))
	}

	return hash
}