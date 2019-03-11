package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// 发送value,结束后关闭channel
func Walk(t *tree.Tree, ch chan int) {
	sendValue(t, ch)
	close(ch)
}

// 递归向channel传值
func sendValue(t *tree.Tree, ch chan int) {
	if t != nil {
		sendValue(t.Left, ch)
		ch <- t.Value
		sendValue(t.Right, ch)
	}
}

// 使用写好的Walk函数来确定两个tree对象  是否一样 原理还是判断value值
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for i := range ch1 { // ch1关闭后 for循环自动跳出
		if i != <-ch2 {
			return false
		}
	}
	return true
}

func main() {
	// 打印 tree.New(1)的值
	var ch = make(chan int)
	go Walk(tree.New(1), ch)
	for v := range ch {
		fmt.Println(v)
	}

	//  比较两个tree的value值是否相等
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
