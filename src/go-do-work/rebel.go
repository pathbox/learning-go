package gdw

import (
	"sync/atomic"

	"github.com/eapache/channels"
)

// Rebel is used to define a goroutine pool whose purpose is to execute fire-and-forget jobs.
type RebelPool struct {
	jobQueue   *channels.InfiniteChannel
	limiter    *channels.ResizableChannel
	queueDepth int64
}

func NewRebelPool(size int) *RebelPool {
	jobQueue := channels.NewInfiniteChannel()
	limiter := channels.NewResizableChannel()
	limiter.Resize(channels.BufferCap(size))
	rebel := &RebelPool{
		jobQueue:   jobQueue,
		limiter:    limiter,
		queueDepth: 0,
	}

	go func() {
		jobQueueOut := jobQueue.Out()
		limiterIn := limiter.In()
		limiterOut := limiter.Out()
		for jobs := range jobQueueOut {
			switch jt := jobs.(type) {

			case Job:
				limiterIn <- true
				atomic.AddInt64(&rebel.queueDepth, -1)
				go func(j Job) {
					j.DoWork()
					<-limiterOut
				}(jt)

			case []Job:
				for _, job := range jt {
					limiterIn <- true
					atomic.AddInt64(&rebel.queueDepth, -1)
					go func(j Job) {
						j.DoWork()
						<-limiterOut
					}(job)
				}

			}
		}
	}()

	return rebel
}

func (r *RebelPool) SetPoolSize(size int) {
	r.limiter.Resize(channels.BufferCap(size))
}

func (r *RebelPool) GetPoolSize() int {
	return int(r.limiter.Cap())
}

func (r *RebelPool) GetQueueDepth() int {
	return int(atomic.LoadInt64(&r.queueDepth))
}

func (r *RebelPool) Add(job Job, amount int) {
	atomic.AddInt64(&r.queueDepth, int64(amount))
	jobs := make([]Job, amount)
	for i := 0; i < amount; i++ {
		jobs[i] = job
	}
	r.jobQueue.In() <- jobs
}

func (r *RebelPool) AddOne(job Job) {
	atomic.AddInt64(&r.queueDepth, 1)
	r.jobQueue.In() <- job
}

func (r *RebelPool) Close() {
	r.jobQueue.Close()
	r.limiter.Close()
}
