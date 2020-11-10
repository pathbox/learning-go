package Concurrence

import "fmt"

// --------------------------- Job ---------------------
type Job interface {
	Do()
}

// --------------------------- Worker ---------------------
type Worker struct {
	JobQueue chan Job
}

func NewWorker() Worker {
	return Worker{JobQueue: make(chan Job)}
}
func (w Worker) Run(wq chan chan Job) {
	go func() {
		for {
			wq <- w.JobQueue
			select {
			case job := <-w.JobQueue:
				//fmt.Println("job.do....")
				job.Do()
			}
		}
	}()
}

// --------------------------- WorkerPool ---------------------
type WorkerPool struct {
	workerlen   int
	JobQueue    chan Job      // 从外界接收Job
	WorkerQueue chan chan Job // 内部的WorkerQueue, Worker Pool, 从中得到一个Worker，Job传给Worker进行消费处理
}

func NewWorkerPool(workerlen int) *WorkerPool {
	return &WorkerPool{
		workerlen:   workerlen,
		JobQueue:    make(chan Job),
		WorkerQueue: make(chan chan Job, workerlen),
	}
}
func (wp *WorkerPool) Run() {
	fmt.Println("初始化worker")
	//初始化worker
	for i := 0; i < wp.workerlen; i++ {
		worker := NewWorker()
		worker.Run(wp.WorkerQueue)
	}
	// 循环获取可用的worker,往worker中写job
	go func() {
		for {
			select {
			case job := <-wp.JobQueue:
				fmt.Println("循环获取可用的worker,往worker中写job")
				worker := <-wp.WorkerQueue
				worker <- job // 从中得到一个Worker，Job传给Worker进行消费处理
			}
		}
	}()
}
