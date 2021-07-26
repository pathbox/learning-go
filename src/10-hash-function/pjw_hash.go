package hashutil

import "unsafe"

func PJWHash(str []byte) uint32 {
	uint32len := uint32(unsafe.Sizeof(*(*uint32)(unsafe.Pointer(nil))))
	var BitsInUnsignedInt uint32 = uint32(uint32len*8)
	var ThreeQuarters uint32 = uint32((BitsInUnsignedInt*3)/4)
	var OneEighth uint32 = uint32(BitsInUnsignedInt/8)
	var HighBits uint32 = uint32((0xFFFFFFFF)<<(BitsInUnsignedInt-OneEighth))
	var hash uint32 = 0

	for i := 0; i < len(str); i++ {
		hash = (hash << OneEighth) + uint32(rune(str[i]))

		if test := hash & HighBits; test != 0 {
			hash = ((hash ^ (test >> ThreeQuarters)) & (^HighBits))
		}
	}

	return hash
}