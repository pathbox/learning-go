package goworkqueue

import (
	"sync"
	"time"
)

type Queue struct {
	Jobs    chan interface{} // 生产者角色job 需要处理的job数量
	done    chan bool
	workers chan chan int // 消费者角色worker  处理job的多线程能力， worker的目的是多线程处理 job
	mux     sync.Mutex
}

func NewQueue(size int, workers int, callback func(interface{}, int)) (q *Queue) {
	q = &Queue{}

	q.Jobs = make(chan interface{}, size)
	q.done = make(chan bool)
	q.workers = make(chan chan int, workers)

	for w := 1; w <= workers; w++ {
		q.workers <- q.worker(w, callback)
	}

	close(q.workers) // queue要关闭了，把worker也 close
	return
}

// 关键逻辑： for循环会阻塞在这里，不断消费job，直到q.done表示queue 要关闭了，才会return
func (q *Queue) worker(id int, callback func(interface{}, int)) (done chan int) {

	done = make(chan int)
	go func() {
	work:
		for { // 不断的循环
			select {
			case <-q.done:
				break work
			case j := <-q.Jobs: // 从jobs中取出job，放入callback中处理
				callback(j, id)
			}
		}
		done <- id // 表示这个worker结束了，返回worker id
		close(done)
	}()

	return done // 这个worker执行结束了， 当 queue关闭了，才会返回，这个worker的for操作就不再进行,这个worker的工作也就结束了
}

// Run blocks until the queue is closed
// NewQueue 和 Run 要在不同的goroutinue中
func (q *Queue) Run() {

	// Wait for workers to be halted
	for w := range q.workers {
		<-w
	}

	// Nothing should still be mindlessly adding jobs
	close(q.Jobs)
}

// Drain queue of jobs
func (q *Queue) Drain(callback func(interface{})) {
	for j := range q.Jobs {
		callback(j)
	}
}

// Close the work queue
func (q *Queue) Close() {
	q.mux.Lock()
	close(q.done)
	q.mux.Unlock()
}

// Closed reports if this queue is already closed
func (q *Queue) Closed() bool {
	q.mux.Lock()
	defer q.mux.Unlock()

	select {
	case <-q.done:
		return true
	default:
		return false
	}
}

// Add jobs to the queue as long as it hasn't be closed
func (q *Queue) Add(job interface{}) (ok bool) {
	q.mux.Lock()
	select {
	case <-q.done:
		ok = false
	case q.Jobs <- job:
		ok = true
	}
	q.mux.Unlock()
	return
}

// SleepUntilTimeOrChanActivity (whichever comes first)
func SleepUntilTimeOrChanActivity(t time.Duration, c chan interface{}) {
	select {
	case <-time.After(t):
	case <-c:
	}
}
