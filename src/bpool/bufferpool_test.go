package bpool

import (
	"bytes"
	"testing"
)

func TestBufferPool(t *testing.T) {
	var size int = 4
	bufPool := NewBufferPool(size)

	b := bufPool.Get()
	bufPool.Put(b)

	for i := 0; i < size*2; i++ {
		bufPool.Put(bytes.NewBuffer([]byte{}))
	}
	close(bufPool.c)
	if len(bufPool.c) != size {
		t.Fatalf("bufferpool size invalid: got %v want %v", len(bufPool.c), size)
	}
}
