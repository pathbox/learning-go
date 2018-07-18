package wpool

import (
	"sync"
)

type Pool struct {
	Worker chan *Worker

	size int

	Jobs chan *Worker

	wg sync.WaitGroup
}

type Worker struct {
	f func() error
}

func NewTask(f func() error) *Worker {
	return &Worker{
		f: f,
	}
}

func NewPool(cap int) *Pool {
	return &Pool{
		Worker: make(chan *Worker),
		size:   cap,
		Jobs:   make(chan *Worker),
	}
}

func (t *Worker) Run(wg *sync.WaitGroup) {
	t.f()
	wg.Done() // 任务结束处
}

func (p *Pool) work() {
	for task := range p.Jobs {
		task.Run(&p.wg)
	}
}

func (p *Pool) Run() {
	for i := 0; i < p.size; i++ {
		go p.work() // 开启size个goroutine进行range p.Jobs
	}

	for task := range p.Worker {
		p.wg.Add(1) // 任务开始处
		p.Jobs <- task
	}
	close(p.Jobs)
	p.wg.Wait() // pool 关闭了，p.Jobs关闭了，但等待还未执行完的f()执行完
}

func (p *Pool) Close() {
	close(p.Worker)
}

/*
使用流程：
p := NewPool(50)
p.Run()

ts := NewTask(func(){
	fmt.Println("Hello World!")
	//...
	return nil
} error)

p.Worker <- ts // p.Worker 是外层接收 Worker的入口

在 pool内部还有一个传递， p.Worker => p.Jobs => task.Run()

p.Worker 和 p.Jobs 是一样类型的 chan *Worker

为什么需要两个chan *Worker呢？有缓冲作用？使用一个是否可行

*/
