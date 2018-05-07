package gracehttp

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	GRACEFUL_ENVIRON_KEY    = "IS_GRACEFUL"
	GRACEFUL_ENVIRON_STRING = GRACEFUL_ENVIRON_KEY + "=1"
	GRACEFUL_LISTENER_FD    = 3
)

// HTTP server that supported graceful shutdown or restart
type Server struct {
	httpServer *http.Server
	listener   net.Listener

	isGraceful   bool
	signalChan   chan os.Signal
	shutdownChan chan bool
}

func NewServer(addr string, handler http.Handler, readTimeout, writeTimeout time.Duration) *Server {
	isGraceful := false
	if os.Getenv(GRACEFUL_ENVIRON_KEY) != "" {
		isGraceful = true
	}
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,

			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		isGraceful:   isGraceful,
		signalChan:   make(chan os.Signal),
		shutdownChan: make(chan bool),
	}
}

func (srv *Server) ListenAndServe() error {
	addr := srv.httpServer.Addr

	if addr == "" {
		addr = ":http"
	}

	ln, err := srv.getNetListener(addr)
	if err != nil {
		return err
	}

	srv.listener = ln // 赋值给srv的listener
	return srv.Serve()
}

func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	addr := srv.httpServer.Addr
	if addr == "" {
		addr = ":https"
	}

	config := &tls.Config{}
	if srv.httpServer.TLSConfig != nil {
		*config = *srv.httpServer.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	ln, err := srv.getNetListener(addr)
	if err != nil {
		return err
	}

	srv.listener = tls.NewListener(ln, config)
	return srv.Serve()
}

func (srv *Server) Serve() error {
	go srv.handleSignals()
	err := srv.httpServer.Serve(srv.listener)

	srv.logf("waiting for connections closed.")
	<-srv.shutdownChan
	srv.logf("all connections closed.")
	return err
}

func (srv *Server) logf(format string, args ...interface{}) {
	pids := strconv.Itoa(os.Getpid())
	format = "[pid " + pids + "] " + format

	if srv.httpServer.ErrorLog != nil {
		srv.httpServer.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

func (srv *Server) getTCPListenerFd() (uintptr, error) {
	file, err := srv.listener.(*net.TCPListener).File() // 返回listener 的文件描述符
	if err != nil {
		return 0, err
	}

	return file.Fd(), nil // 返回listener 的文件描述符
}

func (srv *Server) getNetListener(addr string) (net.Listener, error) {
	var ln net.Listener
	var err error

	if srv.isGraceful {
		file := os.NewFile(GRACEFUL_LISTENER_FD, "") // 新建一个文件描述符，fileNo 号指定为 3
		ln, err = net.FileListener(file)             //FileListener returns a copy of the network listener corresponding to the open file f

		if err != nil {
			err = fmt.Errorf("net.Listen error: %v", err)
			return nil, err
		}
	}
	return ln, nil // 返回这个新的 ln
}

func (srv *Server) handleSignals() {
	var sig os.Signal
	signal.Notify(
		srv.signalChan,
		syscall.SIGTERM,
		syscall.SIGUSR2,
	)

	for {
		sig = <-srv.signalChan // 会一直阻塞在这里，等待signal信号

		switch sig {
		case syscall.SIGTERM: // graceful 关闭http server
			srv.logf("received SIGTERM, graceful shutting down HTTP server.")
			srv.shutdownHTTPServer()
		case syscall.SIGUSR2: // graceful restart http server
			srv.logf("received SIGUSR2, graceful restarting HTTP server.")
			if pid, err := srv.startNewProcess(); err != nil {
				srv.logf("start new process failed: %v, continue serving.", err)
			} else {
				srv.logf("start new process successed, the new pid is %d.", pid)
				srv.shutdownHTTPServer() // 关闭掉老的 http server
			}
		default:
		}
	}
}

func (srv *Server) shutdownHTTPServer() {
	if err := srv.httpServer.Shutdown(context.Background()); err != nil { // 使用了原生http包的 Shutdown，就不用自己在用waitGroup的方式来保证老的连接都执行完
		srv.logf("HTTP server shutdown error: %v", err)
	} else {
		srv.logf("HTTP server shutdown success.")
		srv.shutdownChan <- true //告知老的服务关闭成功
	}
}

// start new process to handle HTTP Connection
// Three Step
func (srv *Server) startNewProcess() (uintptr, error) {
	listenerFd, err := srv.getTCPListenerFd() // 1.新建一个listenerFd， 从当前listenerFd复制而来
	if err != nil {
		return 0, fmt.Errorf("failed to get socket file descriptor: %v", err)
	}

	// set graceful restart env flag
	envs := []string{}
	for _, value := range os.Environ() {
		if value != GRACEFUL_ENVIRON_STRING {
			envs = append(envs, value)
		}
	}

	envs = append(envs, GRACEFUL_ENVIRON_STRING)

	execSpec := &syscall.ProcAttr{ // 2. 准备execSpec属性，用于fork
		Env:   envs,
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), listenerFd},
	}

	fork, err := syscall.ForkExec(os.Args[0], os.Args, execSpec) // 执行fork操作，创建子进程 os.Args[0]就是二进制命令本身
	if err != nil {
		return 0, fmt.Errorf("failed to forkexec: %v", err)
	}

	return uintptr(fork), nil
}
