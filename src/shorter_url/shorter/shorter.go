package shorter

import "strings"

var ALPHABET = strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "")

func GetShortUrl(id int64) string {
	indexAry := Encode62(id)
	return GetString62(indexAry)
}

func Encode62(id int64) []int64 {
	indexAry := []int64{}
	base := int64(len(ALPHABET))

	for id > 0 { // i < 0 时,说明已经除尽了,已经到最高位,数字位已经没有了
		remainder := id % base
		indexAry = append(indexAry, remainder)
		id = id / base
	}

	return indexAry
}

//  输出字符串, 长度不一定为6
func GetString62(indexAry []int64) string {
	result := ""

	for _, val := range indexAry {
		result = result + ALPHABET[val]
	}

	return reverseString(result)
}

// 反转字符串
func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// # 进制的转换:

// # 1 取余操作%,得到最后一位,
// # 2 之后整除/操作,过滤调上一步已经得到的一位
// # 3 重复 1 直到 结果< 0 说明位数已经操作完了
