package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func main() {
	demo1()
	fmt.Println()
	demo2()
}

// 从stdout管道中获取输出值
func demo1() {
	useBufferIo := false
	fmt.Println("Run command `echo -n \"My first command from golang.\"`: ")
	// cmd0 := exec.Command("echo", "-n", "My first command from golang.") // 注册运行的命令
	cmd0 := exec.Command("ps", "aux")
	stdout0, err := cmd0.StdoutPipe() // 得到一个 Stdout pipe 管道，赋值给stdout0
	if err != nil {
		fmt.Printf("Error: Can not obtain the stdout pipe for command No.0: %s\n", err)
		return
	}
	if err := cmd0.Start(); err != nil { // 执行命令, 命令的输出值会传到stdout0管道中
		fmt.Printf("Error: The command No.0 can not be startup: %s\n", err)
		return
	}

	if !useBufferIo {
		var outputBuf0 bytes.Buffer // 定义bytes.Buffer类型 变量
		for {                       // 自己使用for循环方法，构造简单的缓冲循环读取的方式，避免读取大量bytes的时候，消耗光内存
			tempOutput := make([]byte, 5)      // 临时存储slice
			n, err := stdout0.Read(tempOutput) //stdout0 管道中的输出值传给tempOutput
			// fmt.Println(n)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					fmt.Printf("Error: Can not read data from the pipe: %s\n", err)
					return
				}
			}
			if n > 0 {
				fmt.Println("tempOutput: ", tempOutput[:n])
				outputBuf0.Write(tempOutput[:n]) // tempOutput 每次读取到的值，写入到outputBuf0
			}
		}
		fmt.Println("The result: ")
		fmt.Printf("%s\n", outputBuf0.String()) // 输出最后结果
	} else {
		outputBuf0 := bufio.NewReader(stdout0)   // 使用bufio 缓存，直接从stdout0管道中读取字节
		output0, _, err := outputBuf0.ReadLine() // ReadLine方法读取
		if err != nil {
			fmt.Printf("Error: Can not read data from the pipe: %s\n", err)
			return
		}
		fmt.Printf("%s\n", string(output0)) // 输出读取到的值
	}
}

// 把第一个命令的输出内容 通过管道机制传给第二个命令
func demo2() {
	fmt.Println("Run command `ps aux | grep apipe`: ")
	cmd1 := exec.Command("ps", "aux")
	cmd2 := exec.Command("grep", "apipe")
	stdout1, err := cmd1.StdoutPipe() // 命令一 注册 StdoutPipe管道
	if err != nil {
		fmt.Printf("Error: Can not obtain the stdout pipe for command: %s", err)
		return
	}
	if err := cmd1.Start(); err != nil { // 执行命令一
		fmt.Printf("Error: The command can not running: %s\n", err)
		return
	}
	outputBuf1 := bufio.NewReader(stdout1) // New 一个bufio缓冲Reader 4096 的缓冲区
	stdin2, err := cmd2.StdinPipe()        // 命令二 注册 StdinPipe管道
	if err != nil {
		fmt.Printf("Error: Can not obtain the stdin pipe for command: %s\n", err)
		return
	}
	outputBuf1.WriteTo(stdin2) // outputBuf1 写入到stdin2
	var outputBuf2 bytes.Buffer
	cmd2.Stdout = &outputBuf2
	if err := cmd2.Start(); err != nil {
		fmt.Printf("Error: The command can not be startup: %s\n", err)
		return
	}
	err = stdin2.Close() // 关闭stdin2
	if err != nil {
		fmt.Printf("Error: Can not close the stdio pipe: %s\n", err)
		return
	}
	if err := cmd2.Wait(); err != nil { // 命令等待执行完毕
		fmt.Printf("Error: Can not wait for the command: %s\n", err)
		return
	}
	fmt.Printf("%s\n", outputBuf2.Bytes())
}
