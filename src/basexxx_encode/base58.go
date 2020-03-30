package basexxx

import "errors"

// 编解码预设字符集, Base58编码之后字符串依然可以保持字典序
const encodeBase58Map = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var decodeBase58Map [128]byte

var ErrInvalidBase58 = errors.New("invalid base58")

// 为编解码预先初始化好map
func init() {
	for i := 0; i < 128; i++ {
		decodeBase58Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[encodeBase58Map[i]] = byte(i)
	}
}

// Base58 把unit64数值以encodeBase58Map为基值转为字符串
// Base58的基值列表去除了容易混淆的字符，如: 0,O,l等
func Base58(i uint64) string {
	if i < 58 {
		return string(encodeBase58Map[i])
	}

	b := make([]byte, 0, 11)
	for i >= 58 {
		b = append(b, encodeBase58Map[i%58])
		i /= 58
	}
	b = append(b, encodeBase58Map[i])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// ParseBase58 解析Base58字符串
func ParseBase58(s string) (uint64, error) {
	var id uint64

	for _, v := range s {
		if decodeBase58Map[v] == 0xFF {
			return 0, ErrInvalidBase58
		}
		id = id*58 + uint64(decodeBase58Map[v])
	}

	return id, nil
}

// BytesToBase58 把[]byte类型数据以encodeBase58Map为基值转为字符串
func BytesToBase58(bs []byte) string {
	return BytesToBase(bs, Base58)
}
