func keepAlive(c *websocket.Conn, timeout time.Duration) {
  lastResponse := time.Now()
  c.SetPongHandler(func(msg string)error {
    lastResponse = time.Now()
    return nil
    })

  go func() {
    for {
      err := c.WriteMessage(websocket.PingMessage, []byte("keepalive"))
      if err != nil {
        return
      }
      time.Sleep(timeout/2)
      if (time.Now().Sub(lastResponse) > timeout) {
        c.Close()
        return
      }
    }
    }()
}