package simpool

import (
	"log"
)

type SimplePool struct {
	workers chan *woker
}

func NewSimplePool(maxWorkerSize int) *SimplePool {
	newPool := SimplePool{
		workers: make(chan *worker, maxWorkerSize),
	}

	go newPool.workRoutine()

	return &newPool
}

func (p *SimplePool) AddJob(j job) {
	p.workers <- &worker{
		job:             j,
		maxRetryTimes:   j.RetryTimes(),
		alreadyRunTimes: 0,
	}
}

func (p *SimplePool) notifyWorkDone() {
	log.Println("a woker released")
}

func (p *SimplePool) doWorkAsync(w *worker) {
	go func() {
		w.alreadyRunTimes++
		err := w.job.DoJob()
		if err != nil {

			if w.maxRetryTimes >= 0 {
				if w.alreadyRunTimes > w.maxRetryTimes {
					log.Println("already:%d max:%d", w.alreadyRunTimes, w.maxRetryTimes)
				} else {
					p.workers <- w
					log.Println("an old job added")
				}
			} else {
				p.workers <- w
				log.Println("an old job added")
			}
		}

		p.notifyWorkDone()
	}()
}

func (p *SimplePool) workRoutine() {
	for {
		select {
		case w := <-p.workers:
			p.doWorkAsync(w)
		}
	}

}
