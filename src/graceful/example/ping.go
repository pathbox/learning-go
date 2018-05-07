package main

import (
	"syscall"
	"os/signal"
	"os"

	"github.com/Scalingo/go-graceful-restart-example/logger"
	"github.com/Scalingo/go-graceful-restart-example/server"
)

func main() {
	log := logger.New("Server")

	var s *server.Server
	var err error
	if os.Getenv("_GRACEFUL_RESTART") == "true" {
		s, err = server.NewFromFD(log, 3) // 传入的文件描述符指定为3，通过child文件描述符graceful启动server
	} else {
		s, err = server.New(log, 12345) // 普通的启动server
	}
	log.Println("Listen on", s.Addr())

	go s.StartAcceptLoop() // 用一个新的goroutine处理连接读写

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGTERM) // 定义三种类型的notify
	for sig := range signals { // 阻塞，直到接收到signal。 用for{}+select 的处理方式也是可以的
		if sig == syscall.SIGTERM {
			// Stop accepting new connections
			s.Stop()
			err := s.WaitWithTimeout(10 * time.Second)
			if err == server.WaitTimeoutError { // 10s 还没有关闭完
				log.Printf("Timeout when stopping server, %d active connections will be cut.\n", s.ConnectionsCounter())
				os.Exit(-127)
		}
		// Then the program exists
			log.Println("Server shutdown successful")
			os.Exit(0) // 强制离开
	} else if sig == syscall.SIGHUP{ // 平滑重启
		// Stop accepting requests
		s.Stop() // 设置 deadline
		// Get socket file descriptor to pass it to fork
		listenerFD, err := s.listenerFD() // 取到listener服务的文件描述符指针
		if err != nil {
			log.Fatalln("Fail to get socket file descriptor:", err)
		}
		// Set a flag for the new process start process
		os.Setenv("_GRACEFUL_RESTART", "true") // 设置环境变量
		execSpec := &syscall.ProcAttr{ // 定义执行的系统调用环境属性
			Env: os.Environ(),
			Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), listenerFD},
		}
		// Fork exec the new version of your server
		fork, err := syscall.ForkExec(os.Args[0], os.Args, execSpec) // os.Args[0] 是当前go 二进制服务执行文件  fork 子连接
		if err != nil {
			log.Fatalln("Fail to fork", err)
		}
		log.Println("SIGHUP received: fork-exec to", fork)
		// Wait for all conections to be finished
		s.Wait() // 等待原有的connection 读写结束
		log.Println(os.Getpid(), "Server gracefully shutdown")

		// Stop the old server, all the connections have been closed and the new one is running
		os.Exit(0) // 关闭parent process server
	}
}
