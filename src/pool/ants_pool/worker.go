package ants

import (
	"time"
)

// Worker is the actual executor who runs the tasks,
// it starts a goroutine that accepts tasks and
// performs function calls.
type Worker struct {
	// pool who owns this worker.
	pool *Pool

	// task is a job should be done.
	task chan f

	// recycleTime will be update when putting a worker back into queue.
	recycleTime time.Time
}

// run starts a goroutine to repeat the process
// that performs the function calls.
func (w *Worker) run() {
	go func() {
		for f := range w.task {
			if f == nil {
				w.pool.decrRunning()
				return
			}
			f()
			w.pool.putWorker(w)
		}
	}()
}
