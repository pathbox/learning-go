package main

import (
	"./protocol"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	netListen, err := net.Listen("tcp", "localhost:6060")
	CheckError(err)

	defer netListen.Close()

	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		//timeouSec :=10
		//conn.
		Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {

	// 缓冲区，存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//接收解包
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)

	buffer := make([]byte, 1024)
	for {
		fmt.Printf("%s", buffer)
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffer = protocol.Depack(append(tmpBuffer, buffer[:n]...))
	}
	defer conn.Close()
}

func reader(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			Log(string(data))
		}
	}
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

// 开一个goroutine 循环监听心跳，当时间超过设定的时候，就断开连接
// heart beating
func HeartBeating(conn net.Conn, readerChannel chan byte, timeout int) {
	select {
	case fk := <-readerChannel:
		Log(conn.RemoteAddr().String(), "receive data string: ", string(fk))
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		//conn.SetReadDeadline(time.Now().Add(time.Duration(5) * time.Second))
		break
	case <-time.After(time.Second * 5):
		Log("It's really weird to get Nothing!!!")
		conn.Close()
	}
}

func handleConnectionHeart(conn net.Conn, timeout int) {
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			LogErr(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
		Data := (byffer[:n])
		messagenager := make(chan byte)
		// 开一个goroutine 循环监听心跳，当时间超过设定的时候，就断开连接
		go HeartBeating(conn, messagenager, timeout)
		//心跳计时
		go HeartBeating(conn, messnager, timeout)
		//检测每次Client是否有数据传来
		go GravelChannel(Data, messnager)
		Log("receive data length:", n)
		Log(conn.RemoteAddr().String(), "receive data string:", string(Data))
	}
}

func GravelChannel(n []byte, mess chan byte) {
	for _, v := range n {
		mess <- v
	}
	close(mess)
}
