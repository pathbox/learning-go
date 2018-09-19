package main

import (
	"fmt"
	"math/rand"

	xring "github.com/arriqaaq/xring"
)

func main() {
	nodes := []string{"a", "b", "c"}
	cnf := &xring.Config{
		VirtualNodes: 300,
		LoadFactor:   2,
	}
	hashRing := xring.NewRing(nodes, cnf)

	keyCount := 1000000
	distribution := make(map[string]int)
	key := make([]byte, 4)
	for i := 0; i < keyCount; i++ {
		rand.Read(key)
		node, err := hashRing.Get(string(key)) //  根据key得到一个node，相当于这个key就存到这个node上
		if err != nil {
			fmt.Println("error: ", err)
			continue
		}
		hashRing.Done(node)
		distribution[node]++
	}
	for node, count := range distribution {
		fmt.Printf("node: %s, key count: %d\n", node, count)
	}
}
