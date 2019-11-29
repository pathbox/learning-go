package main

import (
	"log"
	"time"

	redis "github.com/pascaldekloe/redis"
)

// Redis is a thread-safe client.
var Redis = redis.NewClient("localhost", time.Second, time.Second)

func main() {
	newLen, err := Redis.RPUSHString("demo_list", "foo")
	if err != nil {
		log.Print("demo_list update error: ", err)
		return
	}
	log.Printf("demo_list has %d elements", newLen)
}
