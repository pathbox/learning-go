package main

import (
	"container/heap"
	"container/ring"
	"fmt"
)

func josephus(n, m int) []int {
	var res []int
	ring := ring.New(n)
	ring.Value = 1
	for i, p := 2, ring.Next(); p != ring; i, p = i+1, p.Next() {
		p.Value = i
	}
	h := ring.Prev()
	for h != h.Next() {
		for i := 1; i < m; i++ {
			h = h.Next()
		}
		res = append(res, h.Unlink(1).Value.(int))
	}
	res = append(res, h.Value.(int))
	return res
}

type intHeap []int

func (h intHeap) Len() int           { return len(h) }
func (h intHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h intHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *intHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	fmt.Println(josephus(9, 5))
	h := &intHeap{10, 3, 9, 7, 2, 88, 31, 67}
	heap.Init(h)
	heap.Push(h, 1)
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}
	fmt.Println()
	heap.Push(h, 100)
	fmt.Println(*h)
}
