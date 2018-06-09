package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	input := "foo bar baz\n 1 2 3"
	scanner := bufio.NewScanner(strings.NewReader(input))
	// scanner.Split(bufio.ScanWords)
	for scanner.Scan() { // 默认以 "\n" 分隔
		fmt.Println(scanner.Text())
	}
}
