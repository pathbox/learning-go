

// 异常抛出后进行不断重试

func waitForResponse() error {
  const timeout = 1 * time.Minute
  deadline := time.Now().Add(timeout)
  url := "www.google.com"

  for tries := 0; time.Now().Before(deadline); tries++{
    _, err := http.Get(url)
    if err == nil {
      return nil
    }
    log.Println("Server not responding.Retrying again")
  }
  return fmt.Errorf("Server %s failed to respond after the %s time", url, timeout)
}

