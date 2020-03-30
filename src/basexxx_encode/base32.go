package basexxx

import "errors"

// 编解码预设字符集, Base58与Base62 编码之后字符串依然可以保持字典序
// Base32基于z-base-32
const encodeBase32Map = "ybndrfg8ejkmcpqxot1uwisza345h769"

var decodeBase32Map [128]byte

var ErrInvalidBase32 = errors.New("invalid base32")

// 为编解码预先初始化好map
func init() {
	for i := 0; i < 128; i++ {
		decodeBase32Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[encodeBase32Map[i]] = byte(i)
	}
}

// Base32 把unit64数值以encodeBase32Map为基值转为字符串
func Base32(i uint64) string {
	if i < 32 {
		return string(encodeBase32Map[i])
	}

	b := make([]byte, 0, 12)
	for i >= 32 { // 不断的从encodeBase32Map中取一个字符，加入到b中
		b = append(b, encodeBase32Map[i%32])
		i /= 32
	}
	b = append(b, encodeBase32Map[i])

	// 反转[]byte
	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// ParseBase32 解析Base32字符串
func ParseBase32(s string) (uint64, error) {
	var id uint64

	for _, v := range s {
		if decodeBase32Map[v] == 0xFF {
			return 0, ErrInvalidBase32
		}
		id = id*32 + uint64(decodeBase32Map[v])
	}

	return id, nil
}

// BytesToBase32 把[]byte类型数据以encodeBase32Map为基值转为字符串
func BytesToBase32(bs []byte) string {
	return BytesToBase(bs, Base32)
}
