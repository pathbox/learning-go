package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	fileBasedPipe()
	inMemorySyncPipe()
}

func fileBasedPipe() {
	reader, writer, err := os.Pipe() // Go使用系统函数创建管道，并把它的两端封装成两个*os.File类型的值。有系统级别的管道支持
	if err != nil {
		fmt.Printf("Error: Can not create the named pipe: %s\n", err)
	}
	go func() {
		output := make([]byte, 100)
		n, err := reader.Read(output) // 将bytes存入output
		if err != nil {
			fmt.Printf("Error: Can not read data from the named pipe: %s\n", err)
		}
		fmt.Printf("Read %d byte(s). [file-based pipe]\n", n)

	}()
	time.Sleep(200 * time.Millisecond) // 为了让 上面的 goroutinue执行完
	input := make([]byte, 120)
	for i := 0; i <= 90; i++ {
		input[i] = byte(i)
	}
	n, err := writer.Write(input)
	if err != nil {
		fmt.Printf("Error: Can not write data to the named pipe: %s\n", err)
	}
	fmt.Printf("Written %d byte(s). [file-based pipe]\n", n)

}

//The flow: inpupt -> writer.Writer(input) -> Pipe -> reader.Read(output) -> output
/*
命名管道使用需要并发运行，因为命名管道默认会在其中一端还未就绪的时候阻塞另一端的进程。并且管道是单向的
命名管道可以被多路复用，当有多个输入端同时写入数据的时候，我们需要考虑原子性
*/

func inMemorySyncPipe() {
	reader, writer := io.Pipe() // 存于内存中、有原子性操作保证的管道(内存通道)
	go func() {
		output := make([]byte, 100)
		n, err := reader.Read(output)
		if err != nil {
			fmt.Printf("Error: Can not read data from the named pipe: %s\n", err)
		}
		fmt.Printf("Read %d byte(s). [in-memory pipe]\n", n)
	}()
	input := make([]byte, 26)
	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i)
	}
	n, err := writer.Write(input)
	if err != nil {
		fmt.Printf("Error: Can not write data to the named pipe: %s\n", err)
	}
	fmt.Printf("Written %d byte(s). [in-memory pipe]\n", n)
	time.Sleep(200 * time.Millisecond)
}

/* 命名管道

mkfifo -m 644 myfifo1
tee dst.log < myfifo1 &
cat src.log > myfifo1

这样，src.log的数据就会通过管道到达dst.log

src.lgo -> myfifo1 -> dst.lgo

*/
