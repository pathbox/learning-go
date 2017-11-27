// create cancel channels
cancelChans := make([]chan struct{}, workerCount)

for i := 0; i < workerCount; i++{
  go worker(jobChan, cancelChans[i])
}

for i := 0; i < workerCount; i++{
  close(cancelChans[i])
}