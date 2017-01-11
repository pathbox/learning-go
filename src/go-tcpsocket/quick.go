func handleConn(c net.Conn) {
  defer c.Close()
  for{
    // read from connection
    // write to the connection
  }
}

func main() {
  l, err := net.Listen("tcp", ":8080")
  if err != nil {
    fmt.Println("listen error: ", err)
    return
  }

  for {
    c, err := l.Accept()
    if err != nil {
      fmt.Println("accept error: ", err)
      break
    }

    go handleConn(c)
  }
}

// 众所周知，TCP Socket的连接的建立需要经历客户端和服务端的三次握手的过程。连接建立过程中，服务端是一个标准的Listen + Accept的结构(可参考上面的代码)，而在客户端Go语言使用net.Dial或DialTimeout进行连接建立：
// Dial
conn, err := net.Dial("tcp", "google.com:80")
if err != nil {

}

conn, err := net.DialTimeout("tcp", ":8080", 2 * time.Second)
if err != nil {

}

