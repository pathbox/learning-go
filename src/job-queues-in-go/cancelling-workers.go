// create a context that can be cancelled

ctx, cancel := context.WithCancel(context.Background())

// start the goroutine passing it the context
go worker(ctx, jobChan)

func worker(ctx context.Context, jobChan <-chan Job) {
  for {
    select{
    case <-ctx.Done():
      return
    case job := <- jobChan:
      process(job)
    }
  }
}

// Invoke cancel when the worker needs to be stopped. This *does not* wait
// for the worker to exit.
cancel()

