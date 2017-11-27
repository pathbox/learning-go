// TryEnqueue tries to enqueue a job to the given job channel. Returns true if
// the operation was successful, and false if enqueuing would not have been
// possible without blocking. Job is not enqueued in the latter case.

func TryEnqueue(job Job, jobChan <- Job) bool {
  select {
  case jobChan <- job:
    return true
  default:
    return false
  }
}

if !TryEnqueue(job, chan) {
    http.Error(w, "max capacity reached", 503)
    return
}

close(jobChan)