package main

import (
	"flag"
	"fmt"
)

var (
	GitTag    = "2000.01.01.release"
	BuildTime = "2000-01-01T00:00:00+0800"
)

func main() {
	version := flag.Bool("v", false, "version")
	flag.Parse()

	if *version {
		fmt.Println("Git Tag: " + GitTag)
		fmt.Println("Build Time: " + BuildTime)
	}
}

//设置编译参数-ldflags "-w -s"。
// 其中-w为去掉调试信息（无法使用gdb调试），-s为去掉符号表（暂未清楚具体作用）

// go build -ldflags "-X main.VERSION=1.0.0 -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'"
