package basexxx

import "errors"

const encodeBase62Map = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var decodeBase62Map [128]byte

var ErrInvalidBase62 = errors.New("invalid base62")

// 为编解码预先初始化好map
func init() {
	for i := 0; i < 128; i++ {
		decodeBase62Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase62Map); i++ {
		decodeBase62Map[encodeBase62Map[i]] = byte(i)
	}
}

// Base62 把unit64数值以encodeBase62Map为基值转为字符串
func Base62(ui uint64) string {
	if ui < 62 {
		return string(encodeBase62Map[ui])
	}

	b := make([]byte, 0, 11)
	for ui >= 62 {
		b = append(b, encodeBase62Map[ui%62])
		ui /= 62
	}
	b = append(b, encodeBase62Map[ui])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// ParseBase62 解析Base62字符串
func ParseBase62(s string) (uint64, error) {
	var id uint64

	for _, v := range s {
		if decodeBase62Map[v] == 0xFF {
			return 0, ErrInvalidBase62
		}
		id = id*62 + uint64(decodeBase62Map[v])
	}

	return id, nil
}

// BytesToBase62 把[]byte类型数据以encodeBase62Map为基值转为字符串
func BytesToBase62(bs []byte) string {
	return BytesToBase(bs, Base62)
}
