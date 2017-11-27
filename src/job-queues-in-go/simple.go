// simple example

// process the job
func worker(jobChan <-chan Job) {
  for job := range jobChan {
    process(job)
  }
}

// make a channel with a capacity of 100
jobChan := make(chan job, 100)

// start the worker
go worker(jobChan)

// enqueue a job
jobChan <- job