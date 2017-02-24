package main

import (
	"bufio"
	"fmt"
	"io"
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
	inputReader := bufio.NewReader(inputFile)
	for {
		// 我们在一个无限循环中使用 ReadString('\n') 或 ReadBytes('\n') 将文件的内容逐行（行结束符 '\n'）读取出来，Unix和Linux的行结束符是 \n，而Windows的行结束符是 \r\n
		inputString, readerError := inputReader.ReadString('\n')
		// inputString, readerError := inputReader.ReadBytes('\n')
		if readerError == io.EOF {
			return
		}
		fmt.Printf("The input was: %s", inputString)
	}
}
