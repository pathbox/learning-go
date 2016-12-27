package main

import (
	"context"
	"fmt"
	"time"
)

type otherContext struct {
	context.Context
	other string
}

func main() {
	// Background 声称一个context的根，是空的，但不是nil，根context是不能够Cancel的
	bkc := context.Background()
	// 生成一个Child Context ，带有Cancel
	c1, cancel := context.WithCancel(bkc)

	// context 实现了String方法可以被Sprint格式化
	if got, want := fmt.Sprint(c1), "context.Background.WithCancel"; got != want {
		fmt.Errorf("c1.String() = %q want %q", got, want)
	}
	// 所有包含Context接口的数据类型都能够当做Context接口类型使用
	o := otherContext{c1, "wonder"}
	c2, _ := context.WithCancel(o)
	contexts := []context.Context{c1, o, c2}
	/*
	   关系图
	       bkc
	       /
	     c1(o)
	     /
	    c2
	*/
	for i, c := range contexts {
		//Done返回值是一个 <-chan struct{}， 返回nil说明这个Context是不能够用关闭的例如：Background生成的Context
		if d := c.Done(); d == nil {
			fmt.Errorf("c[%d].Done() == %v want non-nil", i, d)
		}
		// context关闭的原因，如果没有则返回nil， context 关闭之后会返回指定的内容
		if e := c.Err(); e != nil {
			fmt.Errorf("c[%d].Err() == %v want nil", i, e)
		}
		// c1和o是同一个context,所以前两个打印出来的内容一样
		tm := c.Done()
		fmt.Println(tm)
		select {
		//此处被block住，执行default，因为context并没有关闭
		case x := <-c.Done():
			fmt.Errorf("<-c.Done() == %v want nothing (it should block)", x)
		default:
		}
	}

	// 关闭c1,也就是关闭了o，那么c2也就被关闭了，关闭之后Done所返回的channel 也就关闭了， 关闭之后的channel能够一直取出值
	cancel()
	time.Sleep(100 * time.Millisecond) // let cancelation propagate
	//前面已经Done之后，再次执行Done
	for i, c := range contexts {
		select {
		// 关闭之后c.Done就不再是空
		case <-c.Done():
		default:
			fmt.Errorf("<-c[%d].Done() blocked, but shouldn't have", i)
		}
		if e := c.Err(); e != context.Canceled {
			fmt.Errorf("c[%d].Err() == %v want %v", i, e, context.Canceled)
		}
	}
}
