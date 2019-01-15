package main

import "fmt"

func main() {
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intStreaam := make(chan int)
		go func() {
			defer close(intStreaam)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStreaam <- i:
				}
			}
		}()
		return intStreaam
	}

	multiply := func(
		done <-chan interface{},
		intStream <-chan int,
		multiplier int,
	) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()
		return multipliedStream
	}

	add := func(
		done <-chan interface{},
		intStream <-chan int,
		additive int,
	) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- i + additive:
				}
			}
		}()
		return addedStream
	}

	done := make(chan interface{})
	defer close(done)

	intStreaam := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}

// 从channel中读取数据，进过处理，传给下一个channel， 通过这些channel，连接成了一个pipelines的构造
// 这是一种pipeline的IO并发读写。假设只有一个channel，一个输入要经过add和multiply之后输出，则下一个输入要等待上一个输出这两个步骤都完成才能进行。 如果使用上述例子的construct pipeline的方式，下一个输入并不需要等待上一个输出两个操作都执行完，只要第一个channel有空位置了，就能进入进行操作，通过这种方式形成了并发。
// 当 add 和multiply 操作耗时会比较大时，pipeline的并发模式效果会越明显。让我想起了Rails单核并发模型，有些类似的思想，减少等待，尽可能让CPU空闲，即使是只能用单核
