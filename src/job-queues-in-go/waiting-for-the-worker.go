var wg sync.WaitGroup

func worker(jobChan <-chan Job) {
  defer wg.Done()

  for job := range jobChan{
    process(job)
  }
}

wg.Add(1)

go worker(jobChan)

close(jobChan)

wg.Wait()