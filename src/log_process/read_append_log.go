package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	path := "/Users/pathbox/code/learning-go/src/log_process/access.log"

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	f.Seek(0, 2) // 设置文件指针，从末尾开始读取数据

	bufRead := bufio.NewReader(f) // 使用bufio 缓冲读取
	for {
		line, err := bufRead.ReadBytes('\n')
		if err == io.EOF { // 如果暂时没有数据进入，则sleep一会
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(err)
		}
		line = line[:len(line)-1] // // 去除末尾换行符
		log.Printf("Log Line: %s\n", line)

	}

}
