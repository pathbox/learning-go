package gopool

import (
	"fmt"
	"time"
)

var ErrScheduleTimeout = fmt.Errorf("schedule error: timed out")

type Pool struct {
	sem  chan struct{}
	work chan func() // 将自定义逻辑的func() 注册传入到 work chan
}

func NewPool(size, queue, spawn int) *Pool {
	if spawn <= 0 && queue > 0 {
		panic("dead queue configuration detected")
	}
	if spawn > size {
		panic("spawn > workers")
	}
	p := &Pool{
		sem:  make(chan struct{}, size),
		work: make(chan func(), queue), // worker 实质就是 一个自定义的func()
	}
	for i := 0; i < spawn; i++ {
		p.sem <- struct{}{}
		go p.worker(func() {})
	}
	return p
}

// Schedule schedules task to be executed over pool's workers.
func (p *Pool) Schedule(task func()) {
	p.schedule(task, nil)
}

// ScheduleTimeout schedules task to be executed over pool's workers.
// It returns ErrScheduleTimeout when no free workers met during given timeout.
func (p *Pool) ScheduleTimeout(timeout time.Duration, task func()) error {
	return p.schedule(task, time.After(timeout))
}

// 调度分几种情况
// 1. 执行超时，返回错误
// 2. p.work 没有阻塞，接收到task
// 3. 直接通过p.sem，go p.worker(task)异步执行task
// 这样可以一定程度上解决 p.work chan阻塞的问题，如果阻塞了，就通过多线程的方式，生成新的goroutine执行task，并且 range p.work，减少p.work 阻塞的情况
func (p *Pool) schedule(task func(), timeout <-chan time.Time) error {
	select {
	case <-timeout:
		return ErrScheduleTimeout
	case p.work <- task:
		return nil
	case p.sem <- struct{}{}:
		go p.worker(task)
		return nil
	}
}

func (p *Pool) worker(task func()) { // 执行 参数的task，以及轮询执行 一次p.work中的func()
	defer func() { <-p.sem }()

	task() // 执行 具体的func()

	for task := range p.work { // 从work 中取出 func() 执行他们,这没有实现参数传递，因为参数可能是动态的
		task()
	}

}
