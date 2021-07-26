package main

import (
	"fmt"

	"github.com/dterei/gotsc"
)

const N = 100

func main() {
	tsc := gotsc.TSCOverhead()
	fmt.Println("TSC Overhead:", tsc)

	start := gotsc.BenchStart()
	for i := 0; i < N; i++ {
		// code to evaluate
	}
	end := gotsc.BenchEnd()
	avg := (end - start - tsc) / N

	fmt.Println("Cycles:", avg)
}
