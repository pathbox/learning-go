package net

import (
	"io"
	"syscall"
)

// TCPListener is a TCP network listener. Clients should typically
// use variables of type Listener instead of assuming TCP.
type TCPListener struct {
	fd *netFD
	lc ListenConfig
}

// Accept implements the Accept method in the Listener interface; it
// waits for the next call and returns a generic Conn.
func (l *TCPListener) Accept() (Conn, error) {
	if !l.ok {
		return nil, syscall.EINVAL
	}

	c, err := l.accept()
	if err != nil {
		return nil, &OpError{Op: "accept", Net: l.fd.net, Source: nil, Addr: l.fd.laddr, Err: err}
	}
	return c, nil
}

func (ln *TCPListener) accept() (*TCPConn, error) {
	fd, err := ln.fd.accept() //  得到一个fd，根据这个fd新建Conn
	if err != nil {
		return nil, err
	}
	tc := newTCPConn(fd)
	if ln.lc.KeepAlive >= 0 {
		setKeepAlive(fd, true)
		ka := ln.lc.KeepAlive
		if ln.lc.KeepAlive == 0 {
			ka = defaultTCPKeepAlive
		}
		setKeepAlivePeriod(fd, ka)
	}
	return tc, nil
}

// TCPConn is an implementation of the Conn interface for TCP network
// connections.
type TCPConn struct {
	conn
}

// Conn
type conn struct {
	fd *netFD
}

type conn struct {
	fd *netFD
}

func (c *conn) ok() bool { return c != nil && c.fd != nil }

// Implementation of the Conn interface.
// Read and Write 实际是通过netFD 的Read 和 Write
// Read implements the Conn Read method.
func (c *conn) Read(b []byte) (int, error) {
	if !c.ok() {
		return 0, syscall.EINVAL
	}
	n, err := c.fd.Read(b)
	if err != nil && err != io.EOF {
		err = &OpError{Op: "read", Net: c.fd.net, Source: c.fd.laddr, Addr: c.fd.raddr, Err: err}
	}
	return n, err
}
func (fd *netFD) Read(p []byte) (n int, err error) {
 n, err = fd.pfd.Read(p)
 runtime.KeepAlive(fd)
 return n, wrapSyscallError("read", err)
}

// Implementation of the Conn interface.

// Read implements the Conn Read method.
func (c *conn) Read(b []byte) (int, error) {
 if !c.ok() {
  return 0, syscall.EINVAL
 }
 n, err := c.fd.Read(b)
 if err != nil && err != io.EOF {
  err = &OpError{Op: "read", Net: c.fd.net, Source: c.fd.laddr, Addr: c.fd.raddr, Err: err}
 }
 return n, err
}

func (fd *netFD) Read(p []byte) (n int, err error) {
 n, err = fd.pfd.Read(p)
 runtime.KeepAlive(fd)
 return n, wrapSyscallError("read", err)
}

// Read implements io.Reader.
func (fd *FD) Read(p []byte) (int, error) {
 if err := fd.readLock(); err != nil {
  return 0, err
 }
 defer fd.readUnlock()
 if len(p) == 0 {
  // If the caller wanted a zero byte read, return immediately
  // without trying (but after acquiring the readLock).
  // Otherwise syscall.Read returns 0, nil which looks like
  // io.EOF.
  // TODO(bradfitz): make it wait for readability? (Issue 15735)
  return 0, nil
 }
 if err := fd.pd.prepareRead(fd.isFile); err != nil {
  return 0, err
 }
 if fd.IsStream && len(p) > maxRW {
  p = p[:maxRW]
 }
 for {
  // 尝试从该 socket 读取数据，因为 socket 在被 listener accept 的时候设置成
  // 了非阻塞 I/O，所以这里同样也是直接返回，不管有没有可读的数据
  n, err := syscall.Read(fd.Sysfd, p)
  if err != nil {
   n = 0
   // err == syscall.EAGAIN 表示当前没有期待的 I/O 事件发生，也就是 socket 不可读
   if err == syscall.EAGAIN && fd.pd.pollable() {
    // 如果当前没有发生期待的 I/O 事件，那么 waitRead
    // 会通过 park goroutine 让逻辑 block 在这里
    if err = fd.pd.waitRead(fd.isFile); err == nil {
     continue
    }
   }

   // On MacOS we can see EINTR here if the user
   // pressed ^Z.  See issue #22838.
   if runtime.GOOS == "darwin" && err == syscall.EINTR {
    continue
   }
  }
  err = fd.eofError(n, err)
  return n, err
 }
}

// Write implements the Conn Write method.
func (c *conn) Write(b []byte) (int, error) {
	if !c.ok() {
		return 0, syscall.EINVAL
	}
	n, err := c.fd.Write(b)
	if err != nil {
		err = &OpError{Op: "write", Net: c.fd.net, Source: c.fd.laddr, Addr: c.fd.raddr, Err: err}
	}
	return n, err
}

func (fd *netFD) accept() (netfd *netFD, err error) {
	// 调用 poll.FD 的 Accept 方法接受新的 socket 连接，返回 socket 的 fd
	d, rsa, errcall, err := fd.pfd.Accept()
	if err != nil {
		if errcall != "" {
			err = wrapSyscallError(errcall, err)
		}
		return nil, err
	}
	// 以 socket fd 构造一个新的 netFD，代表这个新的 socket
	if netfd, err = newFD(d, fd.family, fd.sotype, fd.net); err != nil {
		poll.CloseFunc(d)
		return nil, err
	}
	// 调用 netFD 的 init 方法完成初始化
	if err = netfd.init(); err != nil {
		fd.Close()
		return nil, err
	}
	lsa, _ := syscall.Getsockname(netfd.pfd.Sysfd)
	netfd.setAddr(netfd.addrFunc()(lsa), netfd.addrFunc()(rsa))
	return netfd, nil
}

// Accept wraps the accept network call.
func (fd *FD) Accept() (int, syscall.Sockaddr, string, error) {
 if err := fd.readLock(); err != nil {
  return -1, nil, "", err
 }
 defer fd.readUnlock()

 if err := fd.pd.prepareRead(fd.isFile); err != nil {
  return -1, nil, "", err
 }
 for {
  // 使用 linux 系统调用 accept 接收新连接，创建对应的 socket
  s, rsa, errcall, err := accept(fd.Sysfd)
  // 因为 listener fd 在创建的时候已经设置成非阻塞的了，
  // 所以 accept 方法会直接返回，不管有没有新连接到来；如果 err == nil 则表示正常建立新连接，直接返回
  if err == nil {
   return s, rsa, "", err
  }
  // 如果 err != nil，则判断 err == syscall.EAGAIN，符合条件则进入 pollDesc.waitRead 方法
  switch err {
  case syscall.EAGAIN:
   if fd.pd.pollable() {
    // 如果当前没有发生期待的 I/O 事件，那么 waitRead 会通过 park goroutine 让逻辑 block 在这里
    if err = fd.pd.waitRead(fd.isFile); err == nil {
     continue
    }
   }
  case syscall.ECONNABORTED:
   // This means that a socket on the listen
   // queue was closed before we Accept()ed it;
   // it's a silly error, so try again.
   continue
  }
  return -1, nil, errcall, err
 }
}

// 使用 linux 的 accept 系统调用接收新连接并把这个 socket fd 设置成非阻塞 I/O
ns, sa, err := Accept4Func(s, syscall.SOCK_NONBLOCK|syscall.SOCK_CLOEXEC)
// On Linux the accept4 system call was introduced in 2.6.28
// kernel and on FreeBSD it was introduced in 10 kernel. If we
// get an ENOSYS error on both Linux and FreeBSD, or EINVAL
// error on Linux, fall back to using accept.

// Accept4Func is used to hook the accept4 call.
var Accept4Func func(int, int) (int, syscall.Sockaddr, error) = syscall.Accept4

var serverInit sync.Once

func (pd *pollDesc) init(fd *FD) error {
 serverInit.Do(runtime_pollServerInit)
 ctx, errno := runtime_pollOpen(uintptr(fd.Sysfd))
 if errno != 0 {
  if ctx != 0 {
   runtime_pollUnblock(ctx)
   runtime_pollClose(ctx)
  }
  return syscall.Errno(errno)
 }
 pd.runtimeCtx = ctx
 return nil
}
/*
net.Listen
调用 net.Listen 之后，底层会通过 Linux 的系统调用 socket 方法创建一个 fd 分配给 listener，并用以来初始化 listener 的 netFD ，接着调用 netFD 的 listenStream 方法完成对 socket 的 bind&listen 操作以及对 netFD 的初始化（主要是对 netFD 里的 pollDesc 的初始化），调用链是 runtime.runtime_pollServerInit --> runtime.poll_runtime_pollServerInit --> runtime.netpollGenericInit，主要做的事情是：

调用 epollcreate1 创建一个 epoll 实例 epfd，作为整个 runtime 的唯一 event-loop 使用；
调用 runtime.nonblockingPipe 创建一个用于和 epoll 实例通信的管道，这里为什么不用更新且更轻量的 eventfd 呢？我个人猜测是为了兼容更多以及更老的系统版本；
将 netpollBreakRd 通知信号量封装成 epollevent 事件结构体注册进 epoll 实例

调用链：netFD.init() --> poll.FD.Init() --> poll.pollDesc.init()，最终又会走到这里：
然后把这个 socket fd 注册到 listener 的 epoll 实例的事件队列中去，等待 I/O 事件

我们先来看看 Conn.Read 方法是如何实现的，原理其实和 Listener.Accept 是一样的，具体调用链还是首先调用 conn 的 netFD.Read ，然后内部再调用 poll.FD.Read ，最后使用 Linux 的系统调用 read: syscall.Read 完成数据读取：

*/

/*
runtime.netpoll 的核心逻辑是：

根据调用方的入参 delay，设置对应的调用 epollwait 的 timeout 值；
调用 epollwait 等待发生了可读/可写事件的 fd；
循环 epollwait 返回的事件列表，处理对应的事件类型， 组装可运行的 goroutine 链表并返回。
*/

// netpoll checks for ready network connections.
// Returns list of goroutines that become runnable.
// delay < 0: blocks indefinitely
// delay == 0: does not block, just polls
// delay > 0: block for up to that many nanoseconds
func netpoll(delay int64) gList {
 if epfd == -1 {
  return gList{}
 }

 // 根据特定的规则把 delay 值转换为 epollwait 的 timeout 值
 var waitms int32
 if delay < 0 {
  waitms = -1
 } else if delay == 0 {
  waitms = 0
 } else if delay < 1e6 {
  waitms = 1
 } else if delay < 1e15 {
  waitms = int32(delay / 1e6)
 } else {
  // An arbitrary cap on how long to wait for a timer.
  // 1e9 ms == ~11.5 days.
  waitms = 1e9
 }
 var events [128]epollevent
retry:
 // 超时等待就绪的 fd 读写事件
 n := epollwait(epfd, &events[0], int32(len(events)), waitms)
 if n < 0 {
  if n != -_EINTR {
   println("runtime: epollwait on fd", epfd, "failed with", -n)
   throw("runtime: netpoll failed")
  }
  // If a timed sleep was interrupted, just return to
  // recalculate how long we should sleep now.
  if waitms > 0 {
   return gList{}
  }
  goto retry
 }

 // toRun 是一个 g 的链表，存储要恢复的 goroutines，最后返回给调用方
 var toRun gList
 for i := int32(0); i < n; i++ {
  ev := &events[i]
  if ev.events == 0 {
   continue
  }

  // Go scheduler 在调用 findrunnable() 寻找 goroutine 去执行的时候，
  // 在调用 netpoll 之时会检查当前是否有其他线程同步阻塞在 netpoll，
  // 若是，则调用 netpollBreak 来唤醒那个线程，避免它长时间阻塞
  if *(**uintptr)(unsafe.Pointer(&ev.data)) == &netpollBreakRd {
   if ev.events != _EPOLLIN {
    println("runtime: netpoll: break fd ready for", ev.events)
    throw("runtime: netpoll: break fd ready for something unexpected")
   }
   if delay != 0 {
    // netpollBreak could be picked up by a
    // nonblocking poll. Only read the byte
    // if blocking.
    var tmp [16]byte
    read(int32(netpollBreakRd), noescape(unsafe.Pointer(&tmp[0])), int32(len(tmp)))
    atomic.Store(&netpollWakeSig, 0)
   }
   continue
  }

  // 判断发生的事件类型，读类型或者写类型等，然后给 mode 复制相应的值，
    // mode 用来决定从 pollDesc 里的 rg 还是 wg 里取出 goroutine
  var mode int32
  if ev.events&(_EPOLLIN|_EPOLLRDHUP|_EPOLLHUP|_EPOLLERR) != 0 {
   mode += 'r'
  }
  if ev.events&(_EPOLLOUT|_EPOLLHUP|_EPOLLERR) != 0 {
   mode += 'w'
  }
  if mode != 0 {
   // 取出保存在 epollevent 里的 pollDesc
   pd := *(**pollDesc)(unsafe.Pointer(&ev.data))
   pd.everr = false
   if ev.events == _EPOLLERR {
    pd.everr = true
   }
   // 调用 netpollready，传入就绪 fd 的 pollDesc，
   // 把 fd 对应的 goroutine 添加到链表 toRun 中
   netpollready(&toRun, pd, mode)
  }
 }
 return toRun
}

// netpollready 调用 netpollunblock 返回就绪 fd 对应的 goroutine 的抽象数据结构 g
func netpollready(toRun *gList, pd *pollDesc, mode int32) {
 var rg, wg *g
 if mode == 'r' || mode == 'r'+'w' {
  rg = netpollunblock(pd, 'r', true)
 }
 if mode == 'w' || mode == 'r'+'w' {
  wg = netpollunblock(pd, 'w', true)
 }
 if rg != nil {
  toRun.push(rg)
 }
 if wg != nil {
  toRun.push(wg)
 }
}

// netpollunblock 会依据传入的 mode 决定从 pollDesc 的 rg 或者 wg 取出当时 gopark 之时存入的
// goroutine 抽象数据结构 g 并返回
func netpollunblock(pd *pollDesc, mode int32, ioready bool) *g {
 // mode == 'r' 代表当时 gopark 是为了等待读事件，而 mode == 'w' 则代表是等待写事件
 gpp := &pd.rg
 if mode == 'w' {
  gpp = &pd.wg
 }

 for {
  // 取出 gpp 存储的 g
  old := *gpp
  if old == pdReady {
   return nil
  }
  if old == 0 && !ioready {
   // Only set READY for ioready. runtime_pollWait
   // will check for timeout/cancel before waiting.
   return nil
  }
  var new uintptr
  if ioready {
   new = pdReady
  }
  // 重置 pollDesc 的 rg 或者 wg
  if atomic.Casuintptr(gpp, old, new) {
      // 如果该 goroutine 还是必须等待，则返回 nil
   if old == pdWait {
    old = 0
   }
   // 通过万能指针还原成 g 并返回
   return (*g)(unsafe.Pointer(old))
  }
 }
}

// netpollBreak 往通信管道里写入信号去唤醒 epollwait
func netpollBreak() {
 // 通过 CAS 避免重复的唤醒信号被写入管道，
 // 从而减少系统调用并节省一些系统资源
 if atomic.Cas(&netpollWakeSig, 0, 1) {
  for {
   var b byte
   n := write(netpollBreakWr, unsafe.Pointer(&b), 1)
   if n == 1 {
    break
   }
   if n == -_EINTR {
    continue
   }
   if n == -_EAGAIN {
    return
   }
   println("runtime: netpollBreak write failed with", -n)
   throw("runtime: netpollBreak write failed")
  }
 }
}