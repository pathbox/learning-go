package main

import "fmt"

type T int

func IsClosed(ch <-chan T) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func main() {
	c := make(chan T)
	fmt.Println(IsClosed(c)) // false
	close(c)
	fmt.Println(IsClosed(c)) // true
	fmt.Println(IsClosed(c))
}

// The simple way to check the channel if closed
// 你应该在sender的goroutine关闭channel，从而通知receiver(s)(接收者们)已经没有值可以读了。维持这条原则将保证永远不会发生向一个已经关闭的channel发送值或者关闭一个已经关闭的channel

// 如果你因为某种原因从接收端（receiver side）关闭channel或者在多个发送者中的一个关闭channel，那么你应该使用列在Golang panic/recover Use Cases的函数来安全地发送值到channel中 利用recover 防止 panic
