package main

import "time"

func test() {
	x := 100

	go func(a int) {
		a++
		println(a)
	}(x)

	x++
}

func closureTest() {
	y := 200

	go func() {
		y++
		println(y)
	}()

	y++
}

func main() {
	test()
	closureTest()
	time.Sleep(time.Second)
}

// go build -gcflags "-m" -o test

// # _/home/user/code/learning-go/src/optimize-go/closure/test
// ./main.go:8:5: can inline test.func1
// ./main.go:19:5: can inline closureTest.func1
// ./main.go:8:5: func literal escapes to heap
// ./main.go:8:5: func literal escapes to heap
// ./main.go:19:5: func literal escapes to heap
// ./main.go:19:5: func literal escapes to heap
// ./main.go:20:3: &y escapes to heap
// ./main.go:17:7: moved to heap: y
// 闭包引用原环境变量，导致 y 逃逸到堆上，这必然增加了 GC 扫描和回收对象的数量。

/*
 go build -race -o test && ./test
101
==================
WARNING: DATA RACE
Read at 0x00c42008e000 by goroutine 6:
  main.closureTest.func1()
      /home/user/code/learning-go/src/optimize-go/closure/test/main.go:20 +0x3f

Previous write at 0x00c42008e000 by main goroutine:
  main.closureTest()
      /home/user/code/learning-go/src/optimize-go/closure/test/main.go:24 +0xa7
  main.main()
      /home/user/code/learning-go/src/optimize-go/closure/test/main.go:29 +0x34

Goroutine 6 (running) created at:
  main.closureTest()
      /home/user/code/learning-go/src/optimize-go/closure/test/main.go:19 +0x7d
  main.main()
      /home/user/code/learning-go/src/optimize-go/closure/test/main.go:29 +0x34
==================
202
Found 1 data race(s)
*/
