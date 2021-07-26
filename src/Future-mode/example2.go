package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	req := "执行请求"
	g := new(Group)
	// 借助之前文章提及的waitgroup包装函数，不用也行，运行时加入 '-race'
	var wg waitgroup.WaitGroupWrapper
	for i := 0; i < 1000; i++ {
		wg.Wrap(func() {
			j, _ := g.Do(req, NeedSecondTime)
			fmt.Println("NeedSecondTime被调用了=j", j)
		})
	}
	wg.Wait()
}

// 每执行一次,都需要1秒
var counter = 0

func NeedSecondTime() (interface{}, error) {
	time.Sleep(time.Second)
	counter++
	return counter, nil
}

// call is an in-flight or completed Do call
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group represents a class of work and forms a namespace in which
// units of work can be executed with duplicate suppression.
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized,如果要通用，此处string->inteface{}
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {

	g.mu.Lock()
	// 性能要高点，可以创建时初始化
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// 输入参数是否已经存在，存在就等待对应结果，不用计算。
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	// 不存在，则开始执行fn
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	//获取到计算结果
	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
