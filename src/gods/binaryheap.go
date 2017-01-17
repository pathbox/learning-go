// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/emirpasic/gods/trees/binaryheap"
	"github.com/emirpasic/gods/utils"
)

// BinaryHeapExample to demonstrate basic usage of BinaryHeap

func main() {

	// Min-heap
	heap := binaryheap.NewWithIntComparator()
	heap.Push(2)
	fmt.Println(heap)
	heap.Push(3)
	heap.Push(1)
	fmt.Println(heap.Values())
	_, _ = heap.Peek()
	_, _ = heap.Pop()
	_, _ = heap.Pop()
	heap.Push(1)
	heap.Clear()
	heap.Empty()
	heap.Size()

	// Max-heap
	inverseIntComparator := func(a, b interface{}) int {
		return -utils.IntComparator(a, b)
	}
	heap = binaryheap.NewWith(inverseIntComparator)
	heap.Push(2)
	heap.Push(3)
	heap.Push(1)
	fmt.Println(heap.Values())
}
