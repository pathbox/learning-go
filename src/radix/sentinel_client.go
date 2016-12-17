package main

import (
	"fmt"
	"os"
	//  "time"

	"github.com/mediocregopher/radix.v2/sentinel"
)

func errHandler(err error) {
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}

func getSentinel(address string, poolSize int, names ...string) (reply *sentinel.Client, err error) {
	s, err := sentinel.NewClient("tcp", address, poolSize, names...)
	if err != nil {
		return nil, err
	}
	return s, err
}

func main() {
	s, _ := getSentinel("127.0.0.1:9090", 10, "test")
	s, _ := s.GetMaster("test")
	defer s.PutMaster("test", c)
	c.Cmd("SET", "foo", "bar")
}
