package main

import (
	"fmt"
	"os"
	"runtime/trace"
)

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}

	defer trace.Stop()

	for i := 0; i < 100000; i++ {
		fmt.Println("Hello World!")
	}

	// Your program here
}

// 这将使您的程序以二进制格式在文件trace.out中写入事件数据。 然后可以运行go tool trace trace.out。 这将解析跟踪文件，并使用可视化程序打开浏览器。 该命令还将启动服务器，并使用跟踪数据来响应可视化操作。 在浏览器中加载初始页面后，单击“View trace”。 这将加载跟踪查看器，如上面嵌入的那样。

// 1. go run template.go
// 2. go tool trace trace.out
// 3. open 127.0.0.1:49369
