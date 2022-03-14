package zigzag

import "encoding/binary"

func EncodeInt64(n int64) []byte {
	return compressInt64(int64ToZigZag(n))
}

func DecodeInt64(b []byte) int64 {
	return toInt64(decompressInt64(b))
}

func int64ToZigZag(n int64) int64 {
	return (n << 1) ^ (n >> 63)
}

func toInt64(zz int64) int64 {
	return int64(uint64(zz)>>1) ^ -(zz &1)
}

func compressInt64(zz int64) []byte {
	var res []byte
	size := binary.Size(zz)
	for i := 0; i < size; i++ {
		if (zz & ^0x7F) != 0 {
			res = append(res, byte(0x80|(zz&0x7F)))
			zz = int64(uint64(zz) >> 7)
		} else {
			res = append(res, byte(zz&0x7F))
			break
		}
	}
	return res
}

func decompressInt64(zzByte []byte) int64 {
	var res int64
	for i, offset := 0, 0; i < len(zzByte); i, offset = i+1,offset+7 {
		b := zzByte[i]
		if (b & 0x80) == 0x80 {
			res |= int64(b&0x7F) << offset
		} else {
			res |= int64(b) << offset
			break
		}
	}
	return res
}