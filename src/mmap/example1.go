package main

import (
	"fmt"

	mmap "golang.org/x/exp/mmap"
)

func main() {
	at, err := mmap.Open("/Users/xyz/test_mmap_data.txt")
	if err != nil {
		fmt.Println(err.Error())
	}
	buff := make([]byte, 10)
	//读入的长度为slice预设的长度，0是offset。预设长度过长将会用0填充。
	at.ReadAt(buff, 0)
	fmt.Println(string(buff))
	at.Close()
}

// 能实现通过内存映射的方式，从文件读入数据到内存。避免一次拷贝。

// 将test_mmap_data.txt 文件数据读入内存 mmap中，其他进程可以从mmap中读取该文件数据，起到共享该文件数据的作用(也是一种共享内存),减少对该test_mmap_data.txt文件的read IO操作,减少一次拷贝
