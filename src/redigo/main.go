package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func main() {
	rdsConn, gErr := redis.Dial("tcp", "localhost:6379")
	if gErr != nil {
		fmt.Println(gErr)
		return
	}

	defer rdsConn.Close()

	if _, err := rdsConn.Do("SET", "a", "apple"); err != nil {
		fmt.Println(err)
		return
	}

	if reply, err := rdsConn.Do("GET", "a"); err != nil {
		fmt.Println(err)
		return
	} else {
		if replyBytes, ok := reply.([]byte); ok {
			fmt.Println(string(replyBytes))
		} else {
			fmt.Println("Err: get value by string key")
		}
	}
}
