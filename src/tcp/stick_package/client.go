package main

import (
	"encoding/binary"
	"log"
	"net"
)

var HEAD_SIZE = 2

func main() {

}

func SendStringwithTcp(str string) error {
	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
		return err
	}

	head := make([]byte, HEAD_SIZE)
	content := []byte(str)
	headSize := len(content)
	binary.BigEndian.PutUint16(head, uint16(headSize))

	//先写入head部分，再写入body部分
	_, err = conn.Write(head)
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = conn.Write(content)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
