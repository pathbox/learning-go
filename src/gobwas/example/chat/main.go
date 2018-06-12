package main

import (
	"flag"
	"gopool"
	"log"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"

	"net/http"
	_ "net/http/pprof"
)

var (
	addr      = flag.String("listen", ":3333", "address to bind to")
	debug     = flag.String("pprof", "", "address for pprof http")
	workers   = flag.Int("workers", 128, "max workers count")
	queue     = flag.Int("queue", 1, "workers task queue size")
	ioTimeout = flag.Duration("io_timeout", time.Millisecond*100, "i/o operations timeout")
)

func main() {
	flag.Parse()

	if x := *debug; x != "" {
		log.Printf("starting pprof server on %s", x)
		go func() {
			log.Printf("pprof server error: %v", http.ListenAndServe(x, nil))
		}()
	}

	poller, err := netpoll.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	var (
		pool = gopool.NewPool(*workers, *queue, 1)
		chat = NewChat(pool)
		exit = make(chan struct{})
	)

	handle := func(conn net.Conn) {
		safeConn := deadliner{conn, *ioTimeout}

		hs, err := ws.Upgrade(safeConn)
		if err != nil {
			log.Printf("%s: upgrade error: %v", nameConn(conn), err)
			conn.Close()
			return
		}

		log.Printf("%s: established websocket connection: %+v", nameConn(conn), hs)

		user := chat.Register(safeConn)

		desc := netpoll.Must(netpoll.HandleRead(conn))

		poller.Start(desc, func(ev netpoll.Event) {
			if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
				poller.Stop(desc)
				chat.Remove(user)
				return
			}

			pool.Schedule(func() {
				if err := user.Receive(); err != nil {
					poller.Stop(desc)
					chat.Remove(user)
				}
			})
		})
	}

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("websocket is listening on %s", ln.Addr().String())

	acceptDesc := netpoll.Must(netpoll.HandleListener(
		ln, netpoll.EventRead|netpoll.EventOneShot,
	))

	accept := make(chan error, 1)

	poller.Start(acceptDesc, func(e netpoll.Event) {
		err := pool.ScheduleTimeout(time.Millisecond, func() {
			conn, err := ln.Accept()
			if err != nil {
				accept <- err
				return
			}

			accept <- nil
			handle(conn)
		})
		if err == nil {
			err = <-accept
		}
		if err != nil {
			if err != gopool.ErrScheduleTimeout {
				goto cooldown
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				goto cooldown
			}

			log.Fatalf("accept error: %v", err)

		cooldown:
			delay := 5 * time.Millisecond
			log.Printf("accept error: %v; retrying in %s", err, delay)
			time.Sleep(delay)
		}

		poller.Resume(acceptDesc)
	})
	<-exit
}

unc nameConn(conn net.Conn) string {
	return conn.LocalAddr().String() + " > " + conn.RemoteAddr().String()
}

// deadliner is a wrapper around net.Conn that sets read/write deadlines before
// every Read() or Write() call.
type deadliner struct {
	net.Conn
	t time.Duration
}

func (d deadliner) Write(p []byte) (int, error) {
	if err := d.Conn.SetWriteDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Write(p)
}

func (d deadliner) Read(p []byte) (int, error) {
	if err := d.Conn.SetReadDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Read(p)
}