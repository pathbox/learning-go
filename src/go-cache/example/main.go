package main

import (
	"fmt"
	"time"
)

func main() {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 2 seconds
	c := cache.New(5*time.Second, 2*time.Second)

	c.Set("foo", "bar", cache.DefaultExpiration)

	foo, found := c.Get("foo")
	if found {
		fmt.Println(foo)
	}

	time.Sleep(1 * time.Second)

	foo, found = c.Get("foo")
	if found {
		fmt.Println(foo)
	}

	time.Sleep(1 * time.Second)

	foo, found = c.Get("foo")
	if found {
		fmt.Println(foo)
	}

	time.Sleep(1 * time.Second)

	foo, found = c.Get("foo")
	if found {
		fmt.Println(foo)
	}

	time.Sleep(1 * time.Second)

	foo, found = c.Get("foo")
	if found {
		fmt.Println(foo)
	}
}
