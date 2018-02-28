package main

import (
	"fmt"

	"math/rand"

	"runtime"

	"time"
)

func makeBuffer() []byte {

	return make([]byte, rand.Intn(5000000)+5000000)

}

func main() {
	pool := make([][]byte, 20)

	buffer := make(chan []byte, 5)
	var m runtime.MemStats

	makes := 0

	for {
		var b []byte
		select {
		case b = <-buffer:
		default:
			makes += 1
			b = makeBuffer()
		}
		i := rand.Intn(len(pool))
		if pool[i] != nil {
			select {
			case buffer <- pool[i]:
				pool[i] = nil
			default:
			}
		}

		pool[i] = b
		time.Sleep(time.Second)

		bytes := 0
		for i := 0; i < len(pool); i++ {
			if pool[i] != nil {
				bytes += len(pool[i])
			}
		}

		runtime.ReadMemStats(&m)
		fmt.Printf("%d,%d,%d,%d,%d,%d\n", m.HeapSys, bytes, m.HeapAlloc, m.HeapIdle, m.HeapReleased, makes)
	}
}

// For example, using a channel it's possible to keep a separate pool of buffers
//  that are no longer used and use that pool to retrieve a buffer
//  (or make a new one if the channel is empty)

// https://blog.cloudflare.com/recycling-memory-buffers-in-go/
