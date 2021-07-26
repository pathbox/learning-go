package hashutil

func RSHash(str []byte) uint32 {
	var b uint32 = 378551
	var a uint32 = 63689
	var hash uint32 = 0

	for i := 0; i < len(str); i++ {
		hash = hash * a + uint32(rune(str[i]))
		a = a * b
	}

	return hash
}