// create a cancel channel
cancelChan := make(chan struct{})

// start the goroutine passing it the cancel channel
go worker(jobChan, cancelChan)

func worker(jobChan <-chan Job, cancelChan <-chan struct{}) {
    for {
        select {
        case <-cancelChan:
            return

        case job := <-jobChan:
            process(job)
        }
    }
}

// to cancel the worker, close the cancel channel
close(cancelChan)

