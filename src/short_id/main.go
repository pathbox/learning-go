package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
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

// 62个字符, 需要6bit做索引(2 ^ 6 = 64)
var charTable = [...]rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k',
	'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
	'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6',
	'7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
	'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S',
	'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

func ShortenUrl(url string) []string {
	// 取 md5 值的一部分，然后生成短链接
	shortUrlList := make([]string, 0, 4)
	sumData := md5.Sum([]byte(url))
	// 把md5sum分成4份, 每份4个字节
	for i := 0; i < 4; i++ {
		part := sumData[i*4 : i*4+4]
		// 将4字节当作一个整数
		partUint := binary.BigEndian.Uint32(part)
		shortUrlBuffer := &bytes.Buffer{}
		// 将30bit分成6份, 每份5bit
		for j := 0; j < 6; j++ {
			index := partUint % 62
			shortUrlBuffer.WriteRune(charTable[index])
			partUint = partUint >> 5
		}
		shortUrlList = append(shortUrlList, shortUrlBuffer.String())
	}
	return shortUrlList
}

// https://lengzzz.com/note/short-url-algorithm
