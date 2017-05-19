// can't assign requested address 错误解决

// wrong code

for _, route := range routes {
res, err = Request(route.method, ts.URL+route.path)
if err != nil {
  panic(err)
}
res.Body.Close()
}

// The default HTTP client's Transport does not
// attempt to reuse HTTP/1.0 or HTTP/1.1 TCP connections
// ("keep-alive") unless the Body is read to completion and is
// closed.

// 只有当body读取并关闭，而我上面的代码只是关闭了，Body并没有读取。所以导致了client没有reuse TCP connection
// 是client 没有释放 TCP connection 而不是服务端

// right code
for _, route := range routes {
  res, err := Request(route.method, ts.URL+route.path)
  if err != nil {
    panic(err)
  }
  io.Copy(ioutil.Discard, res.Body)
  res.Body.Close()
}