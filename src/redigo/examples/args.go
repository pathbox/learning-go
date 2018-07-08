package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var p1, p2 struct {
	Title  string `redis:"title"`
	Author string `redis:"author"`
	Body   string `redis:"body"`
}

func main() {
	c, err := dial()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	p1.Title = "Example"
	p1.Author = "Gary"
	p1.Body = "Hello"

	// Args is a helper for constructing command arguments from structured values.
	if _, err := c.Do("HMSET", redis.Args{}.Add("id1").AddFlat(&p1)...); err != nil { // key: id1, value: p1
		fmt.Println(err)
		return
	}

	m := map[string]string{
		"title":  "Example2",
		"author": "Steve",
		"body":   "Map",
	}

	if _, err := c.Do("HMSET", redis.Args{}.Add("id2").AddFlat(m)...); err != nil {
		fmt.Println(err)
		return
	}

	for _, id := range []string{"id1", "id2"} {
		v, err := redis.Values(c.Do("HGETALL", id))
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := redis.ScanStruct(v, &p2); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%+v\n", p2)
	}
}

func dial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, nil
	}
	return c, nil
}

//  Do 操作是进行 redis的命令执行，无论是写操作还是读操作，如果想要取回数据，需要使用redis.Values 方法，得到[]interfaces，然后redis.ScanStruct将数据赋值到对应的struct中，之后就可以操作这个struct啦
