package goworkers

import (
	"sync"
	"sync/atomic"
)

const (
	// The size of the buffered queue where jobs are queued up if no
	// workers are available to process the incoming jobs, unless specified
	defaultQSize = 128
	// A comfortable size for the buffered output channel such that chances
	// for a slow receiver to miss updates are minute
	outputChanSize = 100
)

// GoWorkers is a collection of worker goroutines.
//
// All workers will be killed after Stop() is called if their respective job finishes.

type GoWorkers struct {
	numWorkers uint32
	maxWorkers uint32
	numJobs    uint32
	workerQ    chan func()
	bufferedQ  chan func()
	jobQ       chan func()
	stopping   int32
	done       chan struct{}
	// ErrChan is a safe buffered output channel of size 100 on which error
	// returned by a job can be caught, if any. The channel will be closed
	// after Stop() returns. Valid only for SubmitCheckError() and SubmitCheckResult().
	// You must start listening to this channel before submitting jobs so that no
	// updates would be missed. This is comfortably sized at 100 so that chances
	// that a slow receiver missing updates would be minute.
	ErrChan chan error
	// ResultChan is a safe buffered output channel of size 100 on which error
	// and output returned by a job can be caught, if any. The channels will be
	// closed after Stop() returns. Valid only for SubmitCheckResult().
	// You must start listening to this channel before submitting jobs so that no
	// updates would be missed. This is comfortably sized at 100 so that chances
	// that a slow receiver missing updates would be minute.
	ResultChan chan interface{}
}

type Options struct {
	Workers uint32
	QSize   uint32
}

// Accepts optional Options{} argument.
func New(args ...Options) *GoWorkers {
	gw := &GoWorkers{
		workerQ:    make(chan func()),
		jobQ:       make(chan func()),
		ErrChan:    make(chan error, outputChanSize),
		ResultChan: make(chan interface{}, outputChanSize),
		done:       make(chan struct{}),
	}

	gw.bufferedQ = make(chan func(), defaultQSize)
	if len(args) == 1 {
		gw.maxWorkers = args[0].Workers
		if args[0].QSize > defaultQSize {
			gw.bufferedQ = make(chan func(), args[0].QSize)
		}
	}

	go gw.start()

	return gw
}

// JobNum returns number of active jobs
func (gw *GoWorkers) JobNum() uint32 {
	return atomic.LoadUint32(&gw.numJobs)
}

// WorkerNum returns number of active workers
func (gw *GoWorkers) WorkerNum() uint32 {
	return atomic.LoadUint32(&gw.numWorkers)
}

// Submit is a non-blocking call with arg of type `func()`
func (gw *GoWorkers) Submit(job func()) {
	if atomic.LoadInt32(&gw.stopping) == 1 {
		return
	}
	atomic.AddUint32(&gw.numJobs, uint32(1)) // 增1
	gw.jobQ <- func() { job() }
}

// SubmitCheckError is a non-blocking call with arg of type `func() error`
//
// Use this if your job returns 'error'.
// Use ErrChan buffered channel to read error, if any.
func (gw *GoWorkers) SubmitCheckError(job func() error) {
	if atomic.LoadInt32(&gw.stopping) == 1 {
		return
	}
	atomic.AddUint32(&gw.numJobs, uint32(1))
	gw.jobQ <- func() {
		err := job()
		if err != nil {
			select {
			case gw.ErrChan <- err:
			default:
			}
		}
	}
}

// SubmitCheckResult is a non-blocking call with arg of type `func() (interface{}, error)`
//
// Use this if your job returns output and error.
// Use ErrChan buffered channel to read error, if any.
// Use ResultChan buffered channel to read output, if any.
// For a job, either of error or output would be sent if available.
func (gw *GoWorkers) SubmitCheckResult(job func() (interface{}, error)) {
	if atomic.LoadInt32(&gw.stopping) == 1 {
		return
	}
	atomic.AddUint32(&gw.numJobs, uint32(1))
	gw.jobQ <- func() {
		result, err := job()
		if err != nil {
			select {
			case gw.ErrChan <- err:
			default:
			}
		} else {
			select {
			case gw.ResultChan <- result:
			default:
			}
		}
	}
}

// Stop gracefully waits for jobs to finish running.
//
// This is a blocking call and returns when all the active and queued jobs are finished.
func (gw *GoWorkers) Stop() {
	if !atomc.CompareAndSwapInt32(&gw.stopping, 0, 1) {
		return
	}
	if gw.JobNum() != 0 {
		<-gw.done
	}
	close(gw.jobQ)
}

var mx sync.Mutex

func (gw *GoWorkers) spawnWorker() {
	defer mx.Unlock()
	mx.Lock()
	if ((gw.maxWorkers == 0) || (gw.WorkerNum() < gw.maxWorkers)) && (gw.JobNum() > gw.WorkerNum()) {
		go gw.startWorker()
	}
}

func (gw *GoWorkers) start() {
	defer func() { // 把chan 全部关闭
		close(gw.bufferedQ)
		close(gw.workerQ)
		close(gw.ErrChan)
		close(gw.ResultChan)
	}()

	// start 2 workers in advance
	go gw.startWorker()
	go gw.startWorker()

	go func() {
		for {
			select {
			case job, ok := <-gw.bufferedQ: // 从缓存chan中读取传入到workerQ
				if !ok {
					return
				}
				go func() {
					gw.spawnWorker()
					gw.workerQ <- job
				}()
			}
		}
	}()

	for {
		select {
		case job, ok := <-gw.jobQ: // 用ok 来判断job是否取到
			if !ok {
				return
			}
			select {
			case gw.workerQ <- job: // 默认会直接先传入workerQ，
				go gw.spawnWorker()
			default:
				gw.bufferedQ <- job // 如果此时workerQ 满了，则会到bufferedQ 缓存chan中，bufferedQ中的job最终还是会到workerQ
			}
		}
	}
}

// 从 workerQ 中读取函数并执行。消费执行端
func (gw *GoWorkers) startWorker() {
	defer func() {
		atomic.AddUint32(&gw.numWorkers, ^uint32(0))
	}()

	atomic.AddUint32(&gw.numWorkers, 1)
	for job := range gw.workerQ {
		job()
		if (atomic.AddUint32(&gw.numJobs, ^uint32(0)) == 0) && (atomic.LoadInt32(&gw.stopping) == 1) {
			gw.done <- struct{}{}
		}
	}
}
