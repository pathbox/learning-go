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

	title, err := conn.Cmd("HGET", "album:1", "title").Str()
	if err != nil {
		log.Fatal(err)
	}
	// Similarly, get the artist and convert it to a string.
	artist, err := conn.Cmd("HGET", "album:1", "artist").Str()
	if err != nil {
		log.Fatal(err)
	}

	// And the price as a float64...
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
