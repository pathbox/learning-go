package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func main() {
	c, err := dial()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	c.Do("SET", "foo", 1)
	exists, _ := redis.Bool(c.Do("EXISTS", "foo"))
	fmt.Printf("%#v\n", exists)
}

func dial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, nil
	}
	return c, nil
}
