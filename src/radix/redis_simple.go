package main

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/redis"
	"log"
)

func main() {
	// Establish a connection to the Redis server listening on port 6379 of the
	// local machine. 6379 is the default port, so unless you've already
	// changed the Redis configuration file this should work.

	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}

	// Importantly, use defer to ensure the connection is always properly
	// closed before exiting the main() function.
	defer conn.Close()

	// Send our command across the connection. The first parameter to Cmd()
	// is always the name of the Redis command (in this example HMSET),
	// optionally followed by any necessary arguments (in this example the
	// key, followed by the various hash fields and values).

	resp := conn.Cmd("HMSET", "album:1", "title", "Electric Ladyland", "artist", "Jimi Hendrix", "price", 4.95, "likes", 8)
	if resp.Err != nil {
		log.Fatal(err)
	}

	fmt.Println("Electric Ladyland added! ")

}
