package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-connections/proxy"
)

var testBuf = []byte("Buffalo buffalo Buffalo buffalo buffalo buffalo Buffalo buffalo")
var testBufSize = len(testBuf)

type EchoServer interface{
  Run()
  Close()
  LocalAddr() net.Addr
}

type TCPEchoServer struct {
  listener net.Listener
}

func NewEchoServer(proto, addr string) EchoServer {
  var server EchoServer
  if strings.HasPrefix(proto, "tcp"){
    listener, err := net.Listen(proto, addr)
    if err != nil {
      log.Fatal(err)
    }
    server = &TCPEchoServer{listenrt: listener}
  }
  return server
}

func (server *TCPEchoServer) Run(){
  go func(){
    for{
      client, err := server.listener.Accpet()
      if err != nil {
        panic(err)
      }
      go func(client net.Conn){
        if _, err := io.Copy(client, client); err != nil{
          log.Printf("can't echo to the client: %v\n", err.Error())
        }
        client.Close()
      }(client)
    }
  }()
}

func (server *TCPEchoServer) LocalAddr() net.Addr {
  return server.listener.Addr()
}

func (server *TCPEchoServer) Close(){
  server.listener.Close()
}

func TestTCP4Proxy(){
  backend := NewEchoServer("tcp", "127.0.0.1:9000")
  defer backend.Close()
  backend.Run()
  frontendAddr := &net.TCPAddr{IP: net.IPv4(127.0.0.1), Port: 9000}
  proxy, err := proxy.NewProxy(frontendAddr, backend.LocalAddr())
  if err != nil {
    log.Fatal(err)
  }
  testProxy("tcp", proxy)
}
func testProxy(proto string, proxy proxy.Proxy) {
	testProxyAt(proto, proxy, proxy.FrontendAddr().String())
}

func testProxyAt(proto string, proxy proxy.Proxy, addr string) {
  defer proxy.Close()
  go proxy.Run()
  client, err := net.Dial(proto, addr)
  if err != nil {
    log.Fatalf("Can't connect to the proxy: %v", err)
  }
  defer client.Close()
  client.SetDeadline(time.Now().Add(10 * time.Second))
  if _, err = client.Write(testBuf); err != nil {
		log.Fatal(err)
	}
	recvBuf := make([]byte, testBufSize)
	if _, err = client.Read(recvBuf); err != nil {
		log.Fatal(err)
	}
	if !bytes.Equal(testBuf, recvBuf) {
		log.Fatal(fmt.Errorf("Expected [%v] but got [%v]", testBuf, recvBuf))
	}
}

func main() {
  log.Println("start")
  TestTCP4Proxy()
}
