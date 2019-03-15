package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	for {
		line, err := in.ReadString('\n') // 以 回车结尾 读取一次到line
		if err != nil {
			os.Exit(-1)
		}
		line = strings.Replace(line, "吗", "", -1)
		line = strings.Replace(line, "?", "!", -1)
		fmt.Println(line)
	}
}

/*
在吗
在

你好
你好

能听懂汉语吗?
能听懂汉语!
*/
