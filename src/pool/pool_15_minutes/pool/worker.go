package pool

import (
	"log"
	"pjob"
)

type Job struct {
	ID  int
	Job string
}

type Worker struct {
	ID            int
	WorkerChannel chan chan Job
	Channel       chan Job
	End           chan struct{}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerChannel <- w.Channel
			select {
			case job := <-w.Channel:
				// do job
				pjob.DoJob(job.Job, job.ID)
			case <-w.End:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	log.Printf("worker [%d] is stopping", w.ID)
	w.End <- struct{}{}
}
