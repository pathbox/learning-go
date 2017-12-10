// pool 的作用是 Object Pool

package main

import "github.com/gobwas/pool"
import "fmt"

func main() {
	p := pool.New(0, 128, func(n int) interface{} {
		// Create some object with size n.
		slicePool := [128]string{}
		return slicePool
	})
	x, n := p.Get(100) // Returns object with size 128.defer pool.Put(x, n)
	x = "good"

	defer pool.Put(x, n)

	fmt.Println("x", x)

	// Work with x.
}
