package main

import (
	"fmt"
	"github.com/oxtoacart/bpool"
)

var bufpool *bpool.BufferPool

func main() {
	bufpool = bpool.NewBufferPool(48)
	showFunction()
}

func showFunction() error {
	buf := bufpool.Get()

	fmt.Println(buf)

	// buf = []byte("nice")
	bufpool.Put(buf)
	return nil
}

var freeList = make(chan *Buffer, 100)
var serverChan = make(chan *Buffer)

func client() {
	for {
		var b *Buffer
		// Grab a buffer if avaliable; allocate if not
		select {
		case b = <-freeList:
			// Got one; nothing more to do
		default:
			// None free, so allocate a new one
			b = new(Buffer)
		}
		load(b)         // Read next message from the net
		serverChan <- b // Send to server
	}
}

func server() {
	for {
		b := <-serverChan // wait for work
		process(b)
		// reuse buffer if there is room
		select {
		case freeList <- b:
			// buffer on free list; nothing more to do
		default:
			// free list full just carry on
		}
	}
}
