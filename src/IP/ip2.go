package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "baidu.com:80")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer conn.Close()
	fmt.Println(conn.LocalAddr().String())
	fmt.Println(strings.Split(conn.LocalAddr().String(), ":")[0])
}

// 比如得到的结果

//192.168.1.82:53610 // 使用了这个接口来访问baidu.com:80 端口是随机分配的
//192.168.1.82
