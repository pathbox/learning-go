package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/nu7hatch/gouuid"
)

func main() {
	sid, _ := NewShortID()

	fmt.Println("Short id: ", sid)
}

func NewShortID() (string, error) {
	var out string
	base32 := [32]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5"}

	u4, err := uuid.NewV4()

	if err != nil {
		log.Println(err)
		return "", err
	}

	h := md5.New()

	io.WriteString(h, u4.String())

	hexString := fmt.Sprintf("%x", h.Sum(nil))

	intValue, err := strconv.ParseInt(hexString[(31-8):31], 16, 64)
	if err != nil {
		return "", err
	}

	oriole := 0x3FFFFFFF & intValue

	fmt.Println("oriole: ", oriole)

	for j := 0; j < 6; j++ { // j < 6,得到一个6位长度的短id
		val := 0x0000001F & oriole

		fmt.Println("val: ", val)
		out = out + base32[val]
		oriole = oriole >> 5 // 每次获得一个val，将oriole 右移5位，这样得到一个新的oriole
	}
	return out, nil
}
