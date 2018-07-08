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

	c.Do("SADD", "set_with_integers", 4, 5, 6)
	ints, _ := redis.Ints(c.Do("SMEMBERS", "set_with_integers"))

	fmt.Printf("%#v\n", ints)
}

func dial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, nil
	}
	return c, nil
}
