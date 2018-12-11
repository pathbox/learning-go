package gpool

type Pool struct {
	work chan func()
	sem  chan struct{}
}

func New(size int) *Pool {
	return &Pool{
		work: make(chan func()),
		sem:  make(chan struct{}, size),
	}
}

func (p *Pool) Schedule(task func()) error {
	select {
	case p.work <- task: //之后的Schedule走的是这里，因为第一次的时候p.work 是阻塞状态的，没有<-p.work
	case p.sem <- struct{}{}: //第一次Schedule的时候走的是这里
		go p.worker(task) // 执行任务
	}
	return nil
}

func (p *Pool) worker(task func()) {
	defer func() { <-p.sem }()
	for {
		task()
		task = <-p.work
	}
}

// 之后每次Schedule调用代码的执行逻辑是 case p.work <- task   task = <-p.work  task()

// gpool 更像是一个简单的channel队列
