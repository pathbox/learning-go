package main

import (
	"log"
	"net"
)

func main() {
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":9090")
	if err != nil {
		log.Println("dial error: ", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")
}

// 如果传给Dial的Addr是可以立即判断出网络不可达，或者Addr中端口对应的服务没有启动，端口未被监听，Dial会几乎立即返回错误，比如：
