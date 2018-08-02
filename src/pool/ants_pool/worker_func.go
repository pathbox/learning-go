package ants

import (
	"time"
)

// WorkerWithFunc is the actual executor who runs the tasks,
// it starts a goroutine that accepts tasks and
// performs function calls.
type WorkerWithFunc struct {
	// pool who owns this worker.
	pool *PoolWithFunc

	// args is a job should be done.
	args chan interface{}

	// recycleTime will be update when putting a worker back into queue.
	recycleTime time.Time
}

// run starts a goroutine to repeat the process
// that performs the function calls.
func (w *WorkerWithFunc) run() {
	go func() {
		for args := range w.args {
			if args == nil {
				w.pool.decrRunning()
				return
			}
			w.pool.poolFunc(args)
			w.pool.putWorker(w)
		}
	}()
}
