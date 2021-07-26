package hashutil

func BKDRHash(str []byte) uint32 {
	var seed uint32 = 131 /* 31 131 1313 13131 131313 etc.. */
	var hash uint32 = 0

	for i := 0; i < len(str); i++ {
		hash = (hash * seed) + uint32(rune(str[i]))
	}

	return hash
}