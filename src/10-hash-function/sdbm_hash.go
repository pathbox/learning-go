package hashutil

func SDBMHash(str []byte) uint32 {
	var hash uint32 = 0

	for i := 0; i < len(str); i++ {
		hash = uint32(rune(str[i])) + (hash << 6) + (hash << 16) - hash
	}

	return hash
}