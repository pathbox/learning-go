package gdw

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/eapache/channels"
)

// // Worker is used to define a goroutine pool whose results and/or execution are of interest, thus awaitable through WaitGroup
type WorkerPool struct {
	jobQueue   *channels.InfiniteChannel
	limiter    *channels.ResizableChannel
	queueDepth int64
	wgMap      map[string]*sync.WaitGroup
	sync.RWMutex
}

// fast way to build a rand stirng
var src = rand.NewSource(time.Now().UnixNano())

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func generateToken() string {
	b := make([]byte, 20)
	for i, cache, remain := 20-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func NewWorkerPool(size int) *WorkerPool {
	jobQueue := channels.NewInfiniteChannel()
	limiter := channels.NewResizableChannel()
	limiter.Resize(channels.BufferCap(size))
	worker := &WorkerPool{
		jobQueue:   jobQueue,
		limiter:    limiter,
		queueDepth: 0,
		wgMap:      make(map[string]*sync.WaitGroup),
	}
	worker.wgMap["gdw_main_pool"] = &sync.WaitGroup{}

	go func() {
		jobQueueOut := jobQueue.Out()
		limiterIn := limiter.In()
		limiterOut := limiter.Out()
		for jobs := range jobQueueOut {
			switch jt := jobs.(type) {

			case Job:
				limiterIn <- true
				atomic.AddInt64(&worker.queueDepth, -1)
				go func(j Job) {
					j.DoWork()
					<-limiterOut
					worker.RLock()
					worker.wgMap["gdw_main_pool"].Done()
					worker.RUnlock()
				}(jt)

			case []Job:
				for _, job := range jt {
					limiterIn <- true
					atomic.AddInt64(&worker.queueDepth, -1)
					go func(j Job) {
						j.DoWork()
						<-limiterOut
						worker.RLock()
						worker.wgMap["gdw_main_pool"].Done()
						worker.RUnlock()
					}(job)
				}

			case *batchedJob:
				limiterIn <- true
				atomic.AddInt64(&worker.queueDepth, -1)
				go func(bj *batchedJob) {
					bj.batched.DoWork()
					<-limiterOut
					worker.RLock()
					worker.wgMap[bj.name].Done()
					worker.wgMap["gdw_main_pool"].Done()
					worker.RUnlock()
				}(jt)

			case []*batchedJob:
				for _, job := range jt {
					limiterIn <- true
					atomic.AddInt64(&worker.queueDepth, -1)
					go func(bj *batchedJob) {
						bj.batched.DoWork()
						<-limiterOut
						worker.RLock()
						worker.wgMap[bj.name].Done()
						worker.wgMap["gdw_main_pool"].Done()
						worker.RUnlock()
					}(job)
				}

			}
		}
	}()

	return worker
}

func (w *WorkerPool) NewBatch(name string) (*Batch, error) {
	w.RLock()
	_, ok := w.wgMap[name]
	w.RUnlock()
	if ok {
		return nil, fmt.Errorf("Batch named %s already exists.", name)
	}
	w.Lock()
	w.wgMap[name] = &sync.WaitGroup{}
	w.Unlock()
	return &Batch{
		worker: w,
		name:   name,
	}, nil
}

func (w *WorkerPool) NewTempBatch() *Batch {
	token := generateToken()
	w.Lock()
	w.wgMap[token] = &sync.WaitGroup{}
	w.Unlock()
	return &Batch{
		worker: w,
		name:   token,
	}
}

func (w *WorkerPool) LoadBatch(name string) (*Batch, error) {
	w.RLock()
	_, ok := w.wgMap[name]
	w.RUnlock()
	if !ok {
		return nil, fmt.Errorf("No batch named %s exists.", name)
	}
	return &Batch{
		worker: w,
		name:   name,
	}, nil
}

func (w *WorkerPool) SetPoolSize(size int) {
	w.limiter.Resize(channels.BufferCap(size))
}

func (w *WorkerPool) GetPoolSize() int {
	return int(w.limiter.Cap())
}

func (w *WorkerPool) GetQueueDepth() int {
	return int(atomic.LoadInt64(&w.queueDepth))
}

func (w *WorkerPool) Add(job Job, amount int) {
	w.RLock()
	w.add(job, amount, "gdw_main_pool")
}

func (w *WorkerPool) add(job Job, amount int, batch string) {
	w.wgMap["gdw_main_pool"].Add(amount)
	w.RUnlock()
	atomic.AddInt64(&w.queueDepth, int64(amount))
	switch batch {

	case "gdw_main_pool":
		jobs := make([]Job, amount)
		for i := 0; i < amount; i++ {
			jobs[i] = job
		}
		w.jobQueue.In() <- jobs

	default:
		bjs := make([]*batchedJob, amount)
		for i := 0; i < amount; i++ {
			bj := &batchedJob{
				batched: job,
				name:    batch,
			}
			bjs[i] = bj
		}
		w.jobQueue.In() <- bjs

	}
}

func (w *WorkerPool) AddOne(job Job) {
	w.RLock()
	w.addOne(job, "gdw_main_pool")
}

func (w *WorkerPool) addOne(job Job, batch string) {
	w.wgMap["gdw_main_pool"].Add(1)
	w.RUnlock()
	atomic.AddInt64(&w.queueDepth, 1)
	switch batch {

	case "gdw_main_pool":
		w.jobQueue.In() <- job

	default:
		bj := &batchedJob{
			batched: job,
			name:    batch,
		}
		w.jobQueue.In() <- bj

	}
}

func (w *WorkerPool) Wait() {
	w.RLock()
	wg := w.wgMap["gdw_main_pool"]
	w.RUnlock()
	wg.Wait()
}

func (w *WorkerPool) WaitBatch(batch string) error {
	w.RLock()
	wg, ok := w.wgMap[batch]
	w.RUnlock()
	if !ok {
		return fmt.Errorf("No batch named %s exists.", batch)
	}
	wg.Wait()
	return nil
}

func (w *WorkerPool) CleanBatch(batch string) error {
	w.RLock()
	_, ok := w.wgMap[batch]
	w.RUnlock()
	if !ok {
		return fmt.Errorf("No batch named %s exists.", batch)
	}
	w.Lock()
	delete(w.wgMap, batch)
	w.Unlock()
	return nil
}

func (w *WorkerPool) Close() {
	w.jobQueue.Close()
	w.limiter.Close()
}
