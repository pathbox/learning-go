package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var (
	server   *http.Server
	listener net.Listener
	graceful = flag.Bool("graceful", false, "listen on fd open 3 (internal use only)")
)

func sleep(w http.ResponseWriter, r *http.Request) {
	duration, err := time.ParseDuration(r.FormValue("duration"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	time.Sleep(duration)
	w.Write([]byte("Hello World"))
}

func main() {
	flag.Parse()

	http.HandleFunc("/sleep", sleep)
	server = &http.Server{Addr: ":5007"}

	var err error
	if *graceful {
		log.Print("main: Listening to existing file descriptor 3.")
		// cmd.ExtraFiles: If non-nil, entry i becomes file descriptor 3+i.
		// when we put socket FD at the first entry, it will always be 3(0+3)
		f := os.NewFile(3, "")
		listener, err = net.FileListener(f)
	} else {
		log.Print("main: Listening on a new file descriptor.")
		listener, err = net.Listen("tcp", server.Addr)
	}

	if err != nil {
		log.Fatalf("listener error: %v", err)
	}

	go func() {
		// server.Shutdown() stops Serve() immediately, thus server.Serve() should not be in main goroutine
		err = server.Serve(listener)
		log.Printf("server.Serve err: %v\n", err)
	}()
	signalHandler()
	log.Printf("signal end")
}
func reload() error {
	tl, ok := listener.(*net.TCPListener)
	if !ok {
		return errors.New("listener is not tcp listener")
	}
	f, err := tl.File()// 新起一个FD
	if err != nil {
		return err
	}
	args := []string{"-graceful"}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// put socket FD at the first entry
	cmd.ExtraFiles = []*os.File{f}
	return cmd.Start()
}

func signalHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	for {
		sig := <-ch
		log.Printf("signal: %v", sig)

		// timeout context for shutdown
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM: // 如果是这两个命令信号，则服务真正关闭
			// stop
			log.Printf("stop")
			signal.Stop(ch)
			server.Shutdown(ctx)
			log.Printf("graceful shutdown")
			return
		case syscall.SIGUSR2: // 如果是这个命令信号，则服务平滑重启
			log.Printf("reload")
			err := reload()
			if err != nil {
				log.Fatalf("graceful restart error: %v", err)
			}
			server.Shutdown(ctx)
			log.Printf("graceful reload")
			return
		}
	}
}
