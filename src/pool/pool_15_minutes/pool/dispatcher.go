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

	for i < workerCount { // init workers 初始化 workerCount 个 worker
		i++
		log.Println("Starting worker: ", i)
		worker := Worker{
			ID:            i,
			Channel:       make(chan Job), // receive the job
			WorkerChannel: WorkerChannel,  // belongs to the gloabl WorkerChannel
			End:           make(chan struct{}),
		}
		worker.Start()                    // 每个worker都起一个goroutine在后台跑着
		workers = append(workers, worker) // stores worker workers就相当于worker pool
	}
	go startCollect(workers, collector.Job, end)
	return collector // collector的作用是，collector.Job将 input 暴露出去到逻辑层使用，job通过input传到collector

}

func startCollect(workers []Worker, input chan Job, end chan struct{}) {
	for {
		select {
		case <-end:
			for _, w := range workers {
				w.Stop() // stop the worker
			}
			return
		case job := <-input:
			worker := <-WorkerChannel // wait for available Worker channel
			worker <- job             // dispatch work(job) to worker
		}
	}
}
