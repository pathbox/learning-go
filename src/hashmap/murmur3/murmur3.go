package murmur3

import (
	"reflect"
	"unsafe"
)

func Sum32Bytes(key []byte) uint32 {
	return Sum32SeedBytes(key, 0)
}

// Sum32SeedBytes returns a hash from the provided key using the specified seed.
func Sum32SeedBytes(key []byte, seed uint32) uint32 {
	return Sum32Seed(*(*string)((unsafe.Pointer)(&reflect.StringHeader{
		Len:  len(key),
		Data: (*reflect.SliceHeader)(unsafe.Pointer(&key)).Data,
	})), seed)
}

// Sum32 returns a hash from the provided key.
func Sum32(key string) uint32 {
	return Sum32Seed(key, 0)
}

// Sum32Seed returns a hash from the provided key using the specified seed.
func Sum32Seed(key string, seed uint32) uint32 {
	var nblocks = len(key) / 4
	var nbytes = nblocks * 4
	var h1 = seed
	const c1 = 0xcc9e2d51
	const c2 = 0x1b873593
	for i := 0; i < nbytes; i += 4 {
		k1 := uint32(key[i+0]) |
			uint32(key[i+1])<<8 |
			uint32(key[i+2])<<16 |
			uint32(key[i+3])<<24

		k1 *= c1
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= c2
		h1 ^= k1
		h1 = (h1 << 13) | (h1 >> 19)
		h1 = h1*5 + 0xe6546b64
	}
	var k1 uint32
	switch len(key) & 3 {
	case 3:
		k1 ^= uint32(key[nbytes+2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(key[nbytes+1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(key[nbytes+0])
		k1 *= c1
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= c2
		h1 ^= k1
	}
	h1 ^= uint32(len(key))
	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16
	return h1
}
