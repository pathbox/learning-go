package hashutil

func DJBHash(str []byte) uint32 {
	var hash uint32 = 5381

	for i := 0; i < len(str); i++ {
		hash = ((hash << 5) + hash) + uint32(rune(str[i]))
	}

	return hash
}