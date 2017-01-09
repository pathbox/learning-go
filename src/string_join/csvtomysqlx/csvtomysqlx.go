package csvtomysqlx

import (
	"bytes"
	"fmt"
	"strings"
)

var count1 int64

func getNextString1() (string, bool) {
	if count1 < 1000000 {
		count1++
		return "test", true
	}
	return "end", false
}

func Buf() {
	var buffer bytes.Buffer
	for {
		if piece, ok := getNextString1(); ok {
			buffer.WriteString(piece)
		} else {
			break
		}
	}
	fmt.Println("buffer string: ", len(buffer.String())/4)
	fmt.Println(len(buffer.String()))
}

var count2 int64

func getNextString2() (string, bool) {
	if count2 < 1000000 {
		count2++
		return "test", true
	}
	return "end", false
}

func Join() {
	lagerslice := make([]string, 1000)
	for {
		if piece, ok := getNextString2(); ok {
			lagerslice = append(lagerslice, piece)
		} else {
			break
		}
	}
	join := strings.Join(lagerslice, "")
	fmt.Println(len(join) / 4)
	fmt.Println("join string: ", len(lagerslice))
}

// 说明： 两种方式都是拼接100w个字符串，buffer的速度和array的速度是不在一个量级上的。而且数组的方式会带来更多的开销。所以对于超长的字符串拼接，建议使用buffer的方式。当然如果是少量的拼接。还是用fmt.Sprintf&+=的方式。
