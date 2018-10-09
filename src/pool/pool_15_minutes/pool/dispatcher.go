package pool

import (
	"log"
)

var WorkerChannel = make(chan chan Job)

// collect jobs
type Collector struct {
	Job chan Job      // receives jobs to send to workers
	End chan struct{} // stop workers
}

func StartDispatcher(workerCount int) Collector {
	var i int
	var workers []Worker
	input := make(chan Job)    // channel to receive work
	end := make(chan struct{}) // channel to spin down workers
	collector := Collector{Job: input, End: end}

	for i < workerCount {
		i++
		log.Println("Starting worker: ", i)
		worker := Worker{
			ID:            i,
			Channel:       make(chan Job), // receive the job
			WorkerChannel: WorkerChannel,  // belongs to the WorkerChannel
			End:           make(chan struct{}),
		}
		worker.Start()
		workers = append(workers, worker) // stores worker
	}
	go startCollect(workers, input, end)
	return collector

}

func startCollect(workers []Worker, input chan Worker, end chan struct{}) {
	for {
		select {
		case <-end:
			for _, w := range workers {
				w.Stop() // stop the worker
			}
			return
		case job := <-input:
			worker := <-WorkerChannel // wait for available Worker
			worker <- job             // dispatch work(job) to worker
		}
	}
}
