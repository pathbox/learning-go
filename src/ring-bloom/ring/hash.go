package ring

const (
	// 128-bit MurmurHash3 constants
	murmur64c1 uint64 = 0x87c37b91114253d5
	murmur64c2 uint64 = 0x4cf5ad432745937f
	murmur64c3 uint64 = 0x52dce729
	murmur64c4 uint64 = 0x38495ab5
	murmur64c5 uint64 = 0xff51afd7ed558ccd
	murmur64c6 uint64 = 0xc4ceb9fe1a85ec53

	// single byte
	single byte = byte(1)
)

// murmur128 returns two 64-bit outputs of a 128-bit MurmurHash3 hash.
func murmur128(data []byte) (uint64, uint64) {
	var h1, h2, k1, k2 uint64
	length := len(data)
	blocks := length / 16

	for i := 0; i < blocks; i++ {
		k1 = bytesToUint64(data[i*16:])
		k2 = bytesToUint64(data[(i*16)+8:])

		k1 *= murmur64c1
		k1 = (k1 << 31) | (k1 >> (64 - 31))
		k1 *= murmur64c2
		h1 ^= k1

		h1 = (h1 << 27) | (h1 >> (64 - 27))
		h1 += h2
		h1 = h1*5 + murmur64c3

		k2 *= murmur64c2
		k2 = (k2 << 33) | (k2 >> (64 - 33))
		k2 *= murmur64c1
		h2 ^= k2

		h2 = (h2 << 31) | (h2 >> (64 - 31))
		h2 += h1
		h2 = h2*5 + murmur64c4
	}

	tail := blocks * 16
	switch length & 15 {
	case 15:
		k2 ^= uint64(data[tail+14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(data[tail+13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(data[tail+12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(data[tail+11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(data[tail+10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(data[tail+9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(data[tail+8])
		k2 *= murmur64c2
		k2 = (k2 << 33) | (k2 >> (64 - 33))
		k2 *= murmur64c1
		h2 ^= k2
		fallthrough

	case 8:
		k1 ^= uint64(data[tail+7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(data[tail+6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(data[tail+5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(data[tail+4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(data[tail+3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(data[tail+2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(data[tail+1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(data[tail])
		k1 *= murmur64c1
		k1 = (k1 << 31) | (k1 >> (64 - 31))
		k1 *= murmur64c2
		h1 ^= k1
	}

	h1 ^= uint64(length)
	h2 ^= uint64(length)

	h1 += h2
	h2 += h1

	h1 = fmix(h1)
	h2 = fmix(h2)

	h1 += h2
	h2 += h1

	return h1, h2
}

// fmix is the 64-bit MurmurHash3 finalizer to avalanche bits.
func fmix(h uint64) uint64 {
	h ^= h >> 33
	h *= murmur64c5
	h ^= h >> 33
	h *= murmur64c6
	h ^= h >> 33
	return h
}

// bytesToUint64 performs little endian conversion from a byte array to an
// unsigned 64-bit int.
func bytesToUint64(b []byte) uint64 {
	_ = b[7] // memory safety
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

// generateMultihash returns 4 64-bit (2 x 128-bit) MurmurHash3 hashes.
func generateMultiHash(data []byte) [4]uint64 {
	h1, h2 := murmur128(data)
	buff := make([]byte, len(data)+1)
	copy(buff, data)
	buff[len(data)] = single
	h3, h4 := murmur128(buff)
	return [4]uint64{h1, h2, h3, h4}
}

// getRound retrieves the simulated nth round of hashing, fed from 4
// pre-generated hashes.
func getRound(hash [4]uint64, n uint64) uint64 {
	index := 2 + (((n + (n % 2)) % 4) / 2)
	pre := hash[n%2]
	post := n * hash[index]
	return pre + post
}
