package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func main() {
	HandleTCP()
}

func HandleTCP() {
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("start listening on 9999")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) error {
	var (
		BUF_SIZE  = 4096
		HEAD_SIZE = 2
		buffer    = bytes.NewBuffer(make([]byte, 0, BUF_SIZE))
		readBytes = make([]byte, BUF_SIZE)
		isHead    = true
		bodyLen   = 0
	)

	for {
		// 读取数据
		readByteNum, err := conn.Read(readBytes)
		if err != nil {
			log.Fatal(err)
			return err
		}
		buffer.Write(readBytes[0:readByteNum]) // 将读取到的数据放到buffer中

		// 然后处理数据
		for {
			if isHead {
				if buffer.Len() >= HEAD_SIZE { // 说明有body数据
					isHead = false
					head := make([]byte, HEAD_SIZE)
					_, err = buffer.Read(head)
					if err != nil {
						log.Fatal(err)
						return err
					}
					bodyLen = int(binary.BigEndian.Uint16(head)) // 把头部的二进制数据转为int整数，就是body len，这样就知道了body的长度
				} else {
					break // 退出for循环
				}
			}
			if !isHead {
				if buffer.Len() >= bodyLen {
					body := make([]byte, bodyLen)
					_, err = buffer.Read(body[:bodyLen])
					if err != nil {
						log.Fatal(err)
						return err
					}
					fmt.Println("received body: " + string(body[:bodyLen]))
					isHead = true
				} else {
					break
				}
			}
		}
	}
	return nil
}
