package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/xtaci/kcp-go"
)

const portEcho = "127.0.0.1:8081"

func listenEcho() (net.Listener, error) {
	return kcp.Listen(portEcho)
}

func handleEcho(sess *kcp.UDPSession) {
	sess.SetWindowSize(4096, 4096)
	sess.SetACKNoDelay(false)
	// NoDelay options
	// fastest: ikcp_nodelay(kcp, 1, 20, 2, 1)
	// nodelay: 0:disable(default), 1:enable
	// interval: internal update timer interval in millisec, default is 100ms
	// resend: 0:disable fast resend(default), 1:enable fast resend
	// nc: 0:normal congestion control(default), 1:disable congestion control
	sess.SetNoDelay(1, 100, 2, 0)
	for {
		buf := make([]byte, 65536)
		n, err := sess.Read(buf)
		log.Println(n, string(buf[:n]))
		// log.Println(len(buf), string(buf[:n]))
		if err != nil {
			panic(err)
		}
		sess.Write(buf[:n])
	}
}

func echoServer() {
	l, err := listenEcho()
	if err != nil {
		log.Println("listenEcho", err)
		panic(err)
	}

	go func() {
		kcplistener := l.(*kcp.Listener)
		kcplistener.SetReadBuffer(4 * 1024 * 1024)
		kcplistener.SetWriteBuffer(4 * 1024 * 1024)
		kcplistener.SetDSCP(46)
		for {
			s, err := l.Accept()
			if err != nil {
				log.Println("Accept", err)
				return
			}
			s.(*kcp.UDPSession).SetReadBuffer(4 * 1024 * 1024)
			s.(*kcp.UDPSession).SetWriteBuffer(4 * 1024 * 1024)
			go handleEcho(s.(*kcp.UDPSession))
		}
	}()
}

func main() {
	echoServer()
	log.Println("echo server listening")
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-ch
	log.Println("get signal", sig)
}
