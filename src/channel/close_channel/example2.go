package main

import (
	"fmt"
	"time"
)

const (
	fmat = "2006-01-02 15:04:05"
)

func main() {
	c := make(chan int)
	go func() {
		time.Sleep(1 * time.Second)
		c <- 10
		close(c)
	}()

	for {
		select {
		case x, ok := <-c:
			fmt.Printf("%v 通道读取到: x=%v ok=%v\n", time.Now().Format(fmat), x, ok)
			time.Sleep(500 * time.Millisecond)
			if !ok {
				c = nil // 把关闭后的通道赋值为nil，则select读取则会阻塞
			}
		default:
			fmt.Printf("%v 通道没有读取到\n", time.Now().Format(fmat))
			time.Sleep(500 * time.Millisecond)
		}
	}
}

/*
x, ok := <-c 返回的值里第一个x是通道内的值，ok是指通道是否关闭，当通道被关闭后，ok则返回false，因此可以根据这个进行操作。读一个已经关闭的通道为什么会出现false，可以看我之前的 对已经关闭的的chan进行读写，会怎么样？为什么？ 。
当返回的ok为false时，执行c = nil 将通道置为nil，相当于读一个未初始化的通道，则会一直阻塞。至于为什么读一个未初始化的通道会出现阻塞，可以看我的另一篇 对未初始化的的chan进行读写，会怎么样？为什么？ 。select中如果任意某个通道有值可读时，它就会被执行，其他被忽略。则select会跳过这个阻塞case，可以解决不断读已关闭通道的问题。
*/
