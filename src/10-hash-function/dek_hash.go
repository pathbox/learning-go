package hashutil

func DEKHash(str []byte) uint32 {
	var hash uint32 = uint32(len(str))

	for i := 0; i < len(str); i++ {
		hash = ((hash << 5) ^ (hash >> 27)) ^ uint32(rune(str[i]))
	}
	return hash
}