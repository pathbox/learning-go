package main

import (
	"fmt"
	"log"

	"github.com/mediocregopher/radix.v2/redis"
)

func main() {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// Issue a HGET command to retrieve the title for a specific album, and use
	// the Str() helper method to convert the reply to a string.

	title, err := conn.Cmd("HGET", "album:1", "title").Str()
	if err != nil {
		log.Fatal(err)
	}

	artist, err := conn.Cmd("HGET", "album:1", "title").Str()
	if err != nil {
		log.Fatal(err)
	}

	price, err := conn.Cmd("HGET", "album:1", "price").Float64()
	if err != nil {
		log.Fatal(err)
	}

	// And the number of likes as an integer.
	likes, err := conn.Cmd("HGET", "album:1", "likes").Int()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s by %s: £%.2f [%d likes]\n", title, artist, price, likes)
}
