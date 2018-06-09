package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	input := "abcdefghijkl"
	scanner := bufio.NewScanner(strings.NewReader(input))
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		fmt.Printf("%t\t%d\t%s\n", atEOF, len(data), data)
		return 0, nil, nil
	}
	scanner.Split(split)
	buf := make([]byte, 2)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Printf("error: %s\n", scanner.Err())
	}
}

// atEOF 这个参数，它被设计用来通知 split 方法是否没有更多的数据可以读取了，如果到达 EOF，或者出现错误，任何一个发生，scanner 就停止不在读取了，这个标记可以返回错误，scanner.Split() 会返回 false 并停止执行。在最后一行的结尾处没有 \n，所以当方法 ScanLines 找不到新的行字符时，它将简单地返回剩余的字符作为最后一个标记
