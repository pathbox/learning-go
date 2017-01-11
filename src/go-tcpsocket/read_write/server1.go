package main

import (
	"log"
	"net"
)

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// read from the connection
		var buf = make([]byte, 10)
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Println("conn read error: ", err)
			return
		}
		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
}

func main() {
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Println("listen error: ", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("accept error: ", err)
			break
		}
		log.Println("accept a new connection")
		go handleConn(c)
	}
}

// 1、Socket中无数据

// 连接建立后，如果对方未发送数据到socket，接收方(Server)会阻塞在Read操作上，这和前面提到的“模型”原理是一致的。执行该Read操作的goroutine也会被挂起。runtime会监视该socket，直到其有数据才会重新
// 调度该socket对应的Goroutine完成read。由于篇幅原因，这里就不列代码了，例子对应的代码文件：go-tcpsock/read_write下的client1.go和server1.go。

// 2、Socket中有部分数据

// 如果socket中有部分数据，且长度小于一次Read操作所期望读出的数据长度，那么Read将会成功读出这部分数据并返回，而不是等待所有期望数据全部读取后再返回。

// 、Socket中有足够数据

// 如果socket中有数据，且长度大于等于一次Read操作所期望读出的数据长度，那么Read将会成功读出这部分数据并返回。这个情景是最符合我们对Read的期待的了：Read将用Socket中的数据将我们传入的slice填满后返回：n = 10, err = nil。

// client端发送的内容长度为15个字节，Server端Read buffer的长度为10，因此Server Read第一次返回时只会读取10个字节；Socket中还剩余5个字节数据，Server再次Read时会把剩余数据读出
