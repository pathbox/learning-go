package hashutil

func FNVHash(str []byte) uint32 {
	var fnvPrime uint32 = 0x811C9DC5
	var hash uint32 = 0

	for i := 0; i < len(str); i++ {
		hash *= fnvPrime
		hash ^= uint32(rune(str[i]))
	}

	return hash
}