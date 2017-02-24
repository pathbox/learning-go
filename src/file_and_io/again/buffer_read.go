package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	inputFile, inputError := os.Open("input.dat")
	if inputError != nil {
		fmt.Printf("An error occurred on opening the inputfile\n" +
			"Does the file exist?\n" +
			"Have you got acces to it?\n")
		return // exit the function on error
	}
	defer inputFile.Close()
	// 读取器(写入器)缓存读取，并不是一次全部放入内存读取，而是分批读取，这样虽然会慢，但是节省了内存

	// 在很多情况下，文件的内容是不按行划分的，或者干脆就是一个二进制文件。在这种情况下，ReadString()就无法使用了，我们可以使用 bufio.Reader 的 Read()，它只接收一个参数
	buf := make([]byte, 1024) // 缓冲器，大小为1024b
	inputReader := bufio.NewReader(inputFile)
	for {
		n, _ := inputReader.Read(buf) // 变量 n 的值表示读取到的字节数
		fmt.Println(n)
		if n == 0 {
			return
		}
	}
}
