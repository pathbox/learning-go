package main

import(
  "fmt"
  "net"
  "io"
  "net/http"
  "time"
)
// http.ListenAndServe(*http_addr, handle)，如何设置keep-alive的时间？
func ListenAndServe(addr string, handler http.Handler, timeout time.Duration) error {
  server := &http.Server{
    Addr: addr,
    Handler: handler,
    ReadTimeoute: timeout,
    }
    return server.ListernAndServe()
}

func main() {
  addr := "127.0.0.1:8080"
  http.HandleFunc("/", func(http.ResponseWriter, *http.Request){})
  go ListernAndServe(addr, nil, time.Second*10)
  time.Sleep(time.Second)
  started := time.Now()
  remoteAddr, _ := net.ResolveTCPAddr("tcp4", addr)
  conn, err := net.DialTCP("tcp4", nil, remoteAddr)
  if err != nil {
    panic("failed to connect")
  }
  defer conn.Close()
  _, err = conn.Read(make([]byte, 128))
  if err != io.EOF{
    panic("should return EOF")
  }
  fmt.Printf("time escaped=%s, error=%s\n", time.Now().Sub(started), err)
}
