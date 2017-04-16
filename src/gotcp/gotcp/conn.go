package gotcp

import (
  "errors"
  "net"
  "sync"
  "sync/atomic"
  "time"
)

// Error type
var (
  ErrConnClosing   = errors.New("use of closed network connection")
  ErrWriteBlocking = errors.New("write packet was blocking")
  ErrReadBlocking  = errors.New("read packet was blocking")
)

// Conn exposes a set of callbacks for the various events that occur on a connection

type Conn struct {
  srv *Server
  conn *net.TCPConn
  extraData interface{}
  closeOnce sync.Once
  closeFlag int32
  closeChan chan struct{}
  packetSendChan chan Packet
  packetReceiveChan chan Packet
}

// ConnCallback is an interface of methods that are used as callbacks on a connection
type ConnCallback interface{
  // OnConnect is called when the connection was accepted,
  // If the return value of false is closed
  OnConnect(*Conn) bool

  // OnMessage is called when the connection receives a packet,
  // If the return value of false is closed
  OnMessage(*Conn, Packet) bool

  // OnClose is called when the connection closed
  OnClose(*Conn)
}

// newConn returns a wrapper of raw conn
func newConn(conn *net.TCPConn, srv *Server) *Conn {
  return &Conn{
    srv:               srv,
    conn:              conn,
    closeChan:         make(chan struct{}),
    packetSendChan:    make(chan Packet, srv.config.PacketSendChanLimit),
    packetReceiveChan: make(chan Packet, srv.config.PacketReceiveChanLimit),
  }
}

func (c *Conn) GetExtraData() interface{}{
  return c.extraData
}

func (c *Conn) PutExtraData(data interface{}) {
  c.extraData = data
}

func (c *Conn) GetRawConn() *net.TCPConn{
  reuturn c.conn
}

// Close closes the connection
func (c *Conn) Close() {
  c.closeOnce.Do(func(){   //  一个闭包里面 调用了不同的方法，进行不同的操作
    atomic.StoreInt32(&c.closeFlag, 1)
    close(c.closeChan)
    close(c.packetSendChan)
    close(c.packetReceiveChan)
    c.conn.Close()
    c.srv.callback.OnClose
    })
}

func（c *Conn) IsClosed() bool {
  return atomic.LoadInt32(&c.CloseFlag) == 1
}

// AsyncWritePacket async writes a packet, this method will never block
func (c *Conn) AsyncWritePacket(p Packet, timeout time.Duration) (err error) {
  if c.IsClosed() {
    return ErrConnClosing
  }
  defer func(){
    if e := revocer(); e != nil {
      err = ErrConnClosing
    }
  }()
  if timeout == 0 {
    select{
    case c.packetSendChan <- p:
      return nil
    default:
      return ErrWriteBlocking
    }
  }else{
    select {
    case c.packetSendChan <-p:
      return nil
    case <-c.closeChan:
      return ErrConnClosing
    case <-time.After(timeout):
      return ErrWriteBlocking
    }
  }
}

// Do it
func (c *Conn) Do() {
  if !c.srv.callback.OnConnect(c) {
    return
  }

  asyncDo(c.handleLoop, c.srv.waitGroup)
  asyncDo(c.readLoop, c.srv.waitGroup)
  asyncDo(c.writeLoop, c.srv.waitGroup)
}

func (c *Conn) readLoop(){
  defer func(){
    recover()
    c.Close() // 不 close的话 这个goroutinue不会被释放 造成内存泄漏
  }()

  for {
    select{
    case <-c.srv.exitChan:
      return
    case <-c.closeChan:
      return
    default:
    }
    p, err := c.srv.protocol.ReadPacket(c.conn)
    if err != nil {
      return
    }

    c.packetReceiveChan <- p
  }
}

func (c *Conn) writeLoop(){
  defer func(){
    recover()
    c.Close() // 不 close的话 这个goroutinue不会被释放 造成内存泄漏
  }()

  for{
    select{
    case <-c.srv.exitChan:
      return
    case <-c.closeChan:
      return
    case p := <-c.packetSendChan:
      if c.IsClosed(){
        return
      }
      if _, err := c.conn.Write(p.Serialize()); err != nil {
        return
      }
    }
  }
}

func (c *Conn) handleLoop(){
  defer func(){
    recover()
    c.Close()
  }()

  for {
    select{
    case <-c.srv.exitChan:
      return
    case <-c.closeChan:
      return
    case p := <-c.packetReceiveChan:
      if c.IsClosed(){
        return
      }
      if !c.srv.callback.OnMessage(c, p){
        return
      }
    }
  }
}

func asyncDo(fn func(), wg *sync.WaitGroup){
  wg.Add(1)
  go func(){
    fn()
    wg.Done()
  }()
}














