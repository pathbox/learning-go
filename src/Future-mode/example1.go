package main

import (
	"fmt"
	"time"
)

// Call的基本定义，对外部使用者的请求、返回以及异步使用进行封装。
type Call struct {
	Request interface{}
	Reply   interface{}
	Done    chan *Call //用于结果返回时,指向自己的指针
}

func (call *Call) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		// 阻塞情况
	}
}

func main() {
	for i := 0; i < 100; i++ {
		var reply *int
		call := GO(i, reply, nil) //获取到了call，但此时call.Reply还不是运算结果
		//先打印结果还没有计算出来的情况
		fmt.Printf("i=%d,运算前：call.Reply=%v \n", i, call.Reply.(*int))

		result := <-call.Done //等待Done的通知，此时call.Reply发生了变化。
		fmt.Printf("i=%d,运算后：call.Reply=%v,result=%+v \n", i, *(call.Reply.(*int)), *(result.Reply.(*int)))
	}
}

// 供业务调用的异步计算函数封装，用户只需要了解对应参数。
func GO(req int, reply *int, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 10)
	} else {
		if cap(done) == 0 {
			fmt.Println("chan容量为0,无法返回结果,退出此次计算!")
			return nil
		}
	}
	call := &Call{
		Request: req,
		Reply:   reply,
		Done:    done,
	}
	//调用一个可能比较耗时的计算，注意用"go"
	go caculate(call)
	return call
}

//真正的业务处理代码
//简单示意,其实存在读写竞争。run -race 就会出现提示
func caculate(call *Call) {
	//假定运算一次需要耗时1秒
	time.Sleep(time.Second)
	tmp := call.Request.(int) * 5
	call.Reply = &tmp
	call.done()
}
