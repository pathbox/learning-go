package simpool

import (
	"container/list"
	"log"
	"sync"
)

type Job interface {
	DoJob() error // 这两个方法要根据具体业务需求实现

	RetryTimes() int
}

type worker struct {
	job             Job
	maxRetryTimes   int
	alreadyRunTimes int
}

type JobPool struct {
	pendingWorkers   *list.List
	maxWorkerNum     int
	runningWorkerNum int
	mu               sync.Mutex
}

var once sync.Once
var poolIns *JobPool

func Instance() *JobPool {
	once.Do(func() {
		poolIns = NewJobPool(10)
	})

	return poolIns
}

// 返回新实例
func NewJobPool(maxWorkerSize int) *JobPool {
	newPool := JobPool{
		pendingWorkers:   list.New(),
		maxWorkerNum:     maxWorkerSize,
		runningWorkerNum: 0,
	}

	return &newPool
}

func (p *JobPool) AddJob(j Job) {
	w := p.getWorker(j)
	if w != nil {
		w.job = j
		w.maxRetryTimes = j.RetryTimes()
		p.doWorkAsync(w)
	} else {
		log.Println(" A new worker added in pending list")

		w = &worker{
			job:             j,
			maxRetryTimes:   j.GetRetryTimes(),
			alreadyRunTimes: 0,
		}
		p.addPendingWorker(w)
	}
}

func (p *JobPool) addPendingWorker(w *worker) {
	p.mu.Lock()
	defer p.mu.Unlock()

	log.Println("added pending worker")

	p.pendingWorkers.PushBack(w)
}

func (p *JobPool) doWorkAsync(w *worker) {
	go func() {
		w.alreadyRunTimes++
		err := w.job.DoJob()
		if err != nil {
			log.Panicln(err)

			if w.maxRetryTimes >= 0 {
				if w.alreadyRunTimes > w.maxRetryTimes {
					log.Println("already:%d max:%d", w.alreadyRunTimes, w.maxRetryTimes)
				} else {
					p.addPendingWorker(w)
				}
			} else {
				p.addPendingWorker(w)
			}
		}
		p.notifyWorkDone()
	}()
}

func (p *JobPool) getWorker(j Job) *worker {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.runningWorkerNum < p.maxWorkerNum {
		w := worker{
			job:             j,
			maxRetryTimes:   j.RetryTimes(),
			alreadyRunTimes: 0,
		}

		p.runningWorkerNum++

		if p.pendingWorkers.Len() > 0 {
			e := p.pendingWorkers.Front()
			p.pendingWorkers.Remove(e)
			p.pendingWorkers.PushBack(&w)
			return e.Value.(*worker)
		}
		return &w
	}
	return nil
}

func (p *JobPool) notifyWorkDone() {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.runningWorkerNum--

	if p.pendingWorkers.Len() > 0 && p.runningWorkerNum < p.maxWorkerNum {
		p.runningWorkerNum++
		e := p.pendingWorkers.Front()
		p.pendingWorkers.Remove(e)

		p.mu.Unlock()

		p.doWorkAsync(e.Value.(*worker))
		return
	}

}
