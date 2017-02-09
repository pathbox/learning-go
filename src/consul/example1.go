client, err := api.NewClient(&api.Config{Address: "127.0.0.1:8500"})

type LockOptions struct {
    Key              string        // Must be set and have write permissions
    Value            []byte        // Optional, value to associate with the lock
    Session          string        // Optional, created if not specified
    SessionOpts      *SessionEntry // Optional, options to use when creating a session
    SessionName      string        // Optional, defaults to DefaultLockSessionName (ignored if SessionOpts is given)
    SessionTTL       string        // Optional, defaults to DefaultLockSessionTTL (ignored if SessionOpts is given)
    MonitorRetries   int           // Optional, defaults to 0 which means no retries
    MonitorRetryTime time.Duration // Optional, defaults to DefaultMonitorRetryTime
    LockWaitTime     time.Duration // Optional, defaults to DefaultLockWaitTime
    LockTryOnce      bool          // Optional, defaults to false which means try forever
}

opts := &api.LockOptions{
  Key: "webhook_receiver/1",
  Value: []byte("set by sender 1"),
  SessionTTL: "10s",
  SessionOpts: &api.SessionEntry{
    Checks: []string{"check1", "check2"},
    Behavior: "release",
  },
}
lock, err := client.LockOpts(opts)

lock, err := client.LockKey("webhook_receiver/1")

stopCh := make(chan struct{})
lockCh, err := lock.Lock(stopCh)
if err != nil {
  panic(err)
}

cancelCtx, cancelRequest := context.WithCancel(context.Background())
req, _ := http.NewRequest("GET", "https://example.com/webhook", nil)
req = req.WithContext(cancelCtx)
go func() {
  http.DefaultClient.Do(req)
  select{
  case <-cancelCtx.Done():
    log.Println("request cancelled")
  default:
    log.Println("request done")
    err = lock.Unlock()
    if err != nil {
     log.Println("lock already unlocked")
   }
  }
}()
go func() {
  <-lockCh
  cancelRequest()
}()