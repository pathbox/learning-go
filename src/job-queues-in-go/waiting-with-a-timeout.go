// WaitTimeout does a Wait on a sync.WaitGroup object but with a specified
// timeout. Returns true if the wait completed without timing out, false
// otherwise.


func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
  ch := make(chan struct{})

  go func(){  // 专门处理 close操作
    wg.Wait()
    close(ch)
    }()

  select {
  case <-ch:
    return true
  case <- time.After(timeout):
    return false
  }
}

// now use the WaitTimeout instead of wg.Wait()
WaitTimeout(&wg, 5 * time.Second)