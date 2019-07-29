package main

import (
	"fmt"
	"math"
)

var base = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

func base62encode(num int) string {
	baseStr := ""
	for {
		if num <= 0 {
			break
		}

		i := num % 62
		baseStr += base[i]
		num = (num - i) / 62
	}
	return baseStr
}

func base62decode(base62 string) int {
	rs := 0
	l := len(base62) // l 就是encode的时候进行了多少次 %和/ 操作
	f := flip(base)
	for i := 0; i < l; i++{
		rs += f[string(base62[i])] * int(math.Pow(62, float64(i)))
	}
	return rs
}

func flip(s []string) map[string]int {
	f := make(map[string]int)
	for index, value := range s {
		f[value] = index
	}
	return f
}

func main() {
	m := base62encode(12645354)
	fmt.Println(m)
	fmt.Println(base62decode(m))
}