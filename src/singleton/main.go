package main

import (
	"fmt"
	"sync"
	"time"
)

// 单例常用]在固定资源的调用上，比如配置文件的加载之后，代码各处都可以使用。

type Single struct {
	Sone *One
	Stwo *Two
}

type One struct {
	One string
}

type Two struct {
	Two string
}

// 定义锁
var (
	singleton *Single
	lock      sync.Mutex
)

func NewSingle() *Single {
	//这里假设有两个goroutine，一个在第一个条件，一个在第二个条件
	//这就避免了加start函数执行两次的问题，也就是是一个实在的单例
	//尽量不要用defer，因为这虽然是方便，但是会带来不必要的开销
	if singleton == nil {
		lock.Lock()
		if singleton == nil {
			fmt.Println("should enter only")
			singleton = Start()
		}
		lock.Unlock()
	}
	return singleton
}

func Start() *Single {
	one := one()
	two := two()

	single := &Single{
		Sone: one,
		Stwo: two,
	}
	return single
}

func one() *One {
	one := &One{
		One: "one",
	}
	return one
}

func two() *Two {
	two := &Two{
		Two: "two",
	}
	return two
}

func main() {
	//1000个goroutine执行单例，可以看到的结果是单例只初始化了一次.不用在每个新的goroutine中都初始化一次
	for i := 0; i < 1000; i++ {
		go func() {
			single := NewSingle()
			fmt.Println(single)
			fmt.Println(single.Sone.One)
		}()
	}
	//防止main.goroutine退出，导致其他goroutine不能被执行
	time.Sleep(time.Second * 6)
}
