package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func main() {
	c, _ := dial()
	var getScript = redis.NewScript(1, `return redis.call('get',KEYS[1])`)

	reply, _ := getScript.Do(c, "foo")
	fmt.Printf("%v\n", reply) // byte value  49 => 1
	fmt.Printf("%s\n", reply)
}

func dial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, nil
	}
	return c, nil
}
