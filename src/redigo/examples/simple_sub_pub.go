package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func main() {
	redisServerAddr := "0.0.0.0:6379"
	c, _ := redis.Dial("tcp", redisServerAddr)
	defer c.Close()

	psConn := redis.PubSubConn{Conn: c}
	channels := []string{"c1", "c2", "c3"}
	if err := psConn.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		return
	}

	go publish()

	for {
		switch n := psConn.Receive().(type) {
		case error:

			return
		case redis.Message:
			if err := onMessage(n.Channel, n.Data); err != nil {
				fmt.Println(err)
				return
			}
		case redis.Subscription:
			switch n.Count {
			case len(channels):
				// Notify application when all channels are subscribed.

			case 0:
				// Return from the goroutine when all channels are unsubscribed.

				return
			}
		}
	}
}

func publish() {
	c, err := dial()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	c.Do("PUBLISH", "c1", "hello")
	c.Do("PUBLISH", "c2", "world")
	c.Do("PUBLISH", "c3", "goodbye")
}

func dial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, nil
	}
	return c, nil
}

func onMessage(channel string, message []byte) error {
	fmt.Printf("channel: %s, message: %s\n", channel, message)
	return nil
}
