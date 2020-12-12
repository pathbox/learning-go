package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/xtaci/kcp-go"
)

const serverPortEcho = "127.0.0.1:8081"
const N = 100

func dialEcho() (*kcp.UDPSession, error) {
	conn, err := kcp.Dial(serverPortEcho)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	return conn.(*kcp.UDPSession), err
}

func test1() {
	sess, err := dialEcho()
	if err != nil {
		panic(err)
	}
	sess.SetWindowSize(4096, 4096)
	sess.SetWriteDelay(true)
	sess.SetACKNoDelay(false)
	// NoDelay options
	// fastest: ikcp_nodelay(kcp, 1, 20, 2, 1)
	// nodelay: 0:disable(default), 1:enable
	// interval: internal update timer interval in millisec, default is 100ms
	// resend: 0:disable fast resend(default), 1:enable fast resend
	// nc: 0:normal congestion control(default), 1:disable congestion control
	sess.SetNoDelay(1, 100, 2, 0)

	for i := 0; i < N; i++ {
		time.Sleep(1 * time.Second)
		data := time.Now().String()
		sess.Write([]byte(data))
		buf := make([]byte, len(data))
		if n, err := io.ReadFull(sess, buf); err == nil {
			log.Println("got len of(data)", n, string(buf))
			if string(buf[:n]) != data {
				log.Println("不一致", n, len([]byte(data)))
			}
		} else {
			panic(err)
		}

	}
	log.Println("test1 done")
	time.Sleep(1 * time.Second)
	sess.Close()
}

func test2() {
	sess, err := dialEcho()
	if err != nil {
		panic(err)
	}
	sess.SetWindowSize(4096, 4096)
	sess.SetWriteDelay(true)
	sess.SetACKNoDelay(false)
	// NoDelay options
	// fastest: ikcp_nodelay(kcp, 1, 20, 2, 1)
	// nodelay: 0:disable(default), 1:enable
	// interval: internal update timer interval in millisec, default is 100ms
	// resend: 0:disable fast resend(default), 1:enable fast resend
	// nc: 0:normal congestion control(default), 1:disable congestion control
	sess.SetNoDelay(1, 100, 2, 0)

	var buffer bytes.Buffer
	for i := 0; i < 1000; i++ {
		buffer.WriteString(fmt.Sprintf("%5d", i))
	}

	bt := buffer.Bytes()

	for i := 0; i < 1; i++ {
		sess.Write(bt)
		buf := make([]byte, len(bt))
		if n, err := io.ReadFull(sess, buf); err == nil {
			log.Println("got len of(data)", n, buffer.String())
			if string(buf[:n]) != buffer.String() {
				log.Println("不一致", n, len(bt))
			}
		} else {
			panic(err)
		}

	}
	time.Sleep(10 * time.Second)
	sess.Close()
	log.Println("test2 done")
}

func main() {
	// 测试小包
	test1()
	// 测试拆包的情况
	// test2()
}
