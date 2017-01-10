package main

import (
	"fmt"
	"math/rand"
	"time"
)

type binFunc func(int, int) int

func main() {
	rand.Seed(time.Now().Unix())

	fns := []binFunc{
		func(x, y int) int { return x + y },
		func(x, y int) int { return x - y },
		func(x, y int) int { return x * y },
		func(x, y int) int { return x / y },
		func(x, y int) int { return x % y },
	}

	fn := fns[rand.Intn(len(fns))]

	x, y := 12, 5
	fmt.Println(fn(x, y))
}

// go的核心概念：函数是一等公民（first class）
