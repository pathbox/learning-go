package hashutil

func ELFHash(str []byte) uint32 {
	var hash uint32 = 0
	var x uint32 = 0

	for i := 0; i < len(str); i++ {
		hash = (hash << 4) + uint32(rune(str[i]))

		if x = hash & 0xF0000000; x != 0 {
			hash ^= (x >> 24)
		}
		hash &= ^x
	}
	return hash
}