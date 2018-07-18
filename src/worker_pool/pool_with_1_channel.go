package wpool

import "sync"

// 只使用一个 chan *Worker代码示例：

type Pool struct {
	Worker chan *Worker

	size int

	CloseChan chan struct{}

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
		Worker:    make(chan *Worker),
		CloseChan: make(chan struct{}),
		size:      cap,
	}
}

func (t *Worker) Run(wg *sync.WaitGroup) {
	t.f()
	wg.Done() // 任务结束处
}

func (p *Pool) work() {
	for task := range p.Worker {
		p.wg.Add(1)
		task.Run(&p.wg)
	}
}

func (p *Pool) Run() {
	for i := 0; i < p.size; i++ {
		go p.work() // 开启size个goroutine进行range p.Jobs
	}

	for {
		select {
		case <-p.CloseChan:
			close(p.Worker) // 关闭 p.Worker
			p.wg.Wait()     // pool 关闭了，p.Worker关闭了，但等待还未执行完的f()执行完
			return
		}
	}

}

func (p *Pool) Close() {
	p.CloseChan <- struct{}{}
}

// 使用一个channel作为 pool 的队列是可以的，再用一个chan来作为关闭通知，甚至可以不用关闭worker pool的操作，一般情况，开启一个worker pool 就是希望一直能够执行下去
