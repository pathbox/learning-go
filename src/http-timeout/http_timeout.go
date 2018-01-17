server := &http.Server{
	ReadTimeout: 30 * time.Second,
	WriteTimeout: 20 * time.Second,
	Addr: "127.0.0.1:9090",
	Handler: router.Handler,
	MaxHeaderBytes: 1 << 20,
}

log.Fatal(server.ListenAndServe())

if d := c.server.ReadTimeout; d != 0 {
	c.rwc.SetReadDeadline(time.Now().Add(d))
}


if d := c.server.WriteTimeout; d != 0 {
	c.rwc.SetReadDeadline(time.Now().Add(d))
}

func (c *conn) readRequest(ctx context.Context) (w *response, err error) {
	if c.hijacked() {
		return nil, ErrHijacked
	}
	if d := c.server.ReadTimeout; d != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(d))
	}
	if d := c.server.WriteTimeout; d != 0 {
		defer func() {
			c.rwc.SetWriteDeadline(time.Now().Add(d))
		}()
	}
  ……
}

// 但是，当连接是HTTPS的时候，SetWriteDeadline会在Accept之后立即调用(代码)，
// 所以它的时间计算也包括 TLS握手时的写的时间。 讨厌的是， 这就意味着(也只有这种情况)
// WriteTimeout设置的时间也包含读取Headerd到读取body第一个字节这段时间。
// 当你处理不可信的客户端和网络的时候，你应该同时设置读写超时，这样客户端就不会因为读慢或者写慢长久的持有这个连接了


if tlsConn, ok := c.rwc.(*tls.Conn); ok {
		if d := c.server.ReadTimeout; d != 0 {
			c.rwc.SetReadDeadline(time.Now().Add(d))
		}
		if d := c.server.WriteTimeout; d != 0 {
			c.rwc.SetWriteDeadline(time.Now().Add(d))
		}
	}

c := &http.Client{
	Timeout: 15 * time.Second,
}

resp, err := c.Get("https://www.baidu.com")

// net.Dialer.Timeout 限制建立TCP连接的时间
// http.Transport.TLSHandshakeTimeout 限制 TLS握手的时间
// http.Transport.ResponseHeaderTimeout 限制读取response header的时间
// http.Transport.ExpectContinueTimeout 限制client在发送包含 Expect: 100-continue的header到收到继续发送body的response之间的时间等待

c := &http.CLient{
	Transport: &Transport{
		Dial: (&net.Dialer{
			Timeout: 30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10*time.Second,
		ResponseHeaderTimeout: 10*time.Second,
		 ExpectContinueTimeout: 1 * time.Second,
	}
}