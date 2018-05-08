package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var (
	graceful = flag.Bool("graceful", false, "-graceful")
)

// Accepted accepted connection
type Accepted struct {
	conn net.Conn
	err  error
}

func handleConnection(conn net.Conn) {
	conn.Write([]byte("hello"))
	conn.Close()
}

func listenAndServe(ln net.Listener, sig chan os.Signal) {
	accepted := make(chan Accepted, 1)
	go func() {
		for {
			conn, err := ln.Accept()
			accepted <- Accepted{conn, err}
		}
	}()

	for {
		select {
		case a := <-accepted:
			if a.err == nil {
				fmt.Println("handle connection")
				go handleConnection(a.conn)
			}
		case _ = <-sig: // 接收到graceful restart signal信号时，执行 graceful restart
			fmt.Println("gonna fork and run")
			forkAndRun(ln)
			break
		}
	}
}

func gracefulListener() net.Listener {
	newFile := os.NewFile(3, "graceful server") // 1. 新建一个文件描述符，然后复制父进程环境信息，并指定给这个文件描述符 得到一个新的监听文件描述符
	ln, err := net.FileListener(newFile)
	if err != nil {
		fmt.Println(err)
	}

	return ln
}

func firstBootListener() net.Listener {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}

	return ln
}

func forkAndRun(ln net.Listener) {
	l := ln.(*net.TCPListener)
	newFile, _ := l.File() // 取出新的监听文件描述符
	fmt.Println(newFile.Fd())

	cmd := exec.Command(os.Args[0], "-graceful")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.ExtraFiles = []*os.File{newFile} // 2. 将新的监听文件描述符配置到cmd属性中，让新的cmd fork命令使用这个新的监听文件描述符
	cmd.Run()                            // cmd 执行fork
}

func main() {
	flag.Parse()
	fmt.Printf("given args: %t, pid: %d\n", *graceful, os.Getpid())
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1) // 定义signal 通知信号

	var ln net.Listener
	if *graceful {
		ln = gracefulListener()
	} else {
		ln = firstBootListener()
	}

	listenAndServe(ln, c)
}
