package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gansidui/gotcp"
	"github.com/gansidui/gotcp/examples/telnet"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// creates a tcp listener

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":23")
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}

	srv := gotcp.NewServer(config, &telnet.TelnetCallback{}, &telnet.TelnetProtocol{})

	go srv.Start(listener, time.Second)
	fmt.Println("listening:", listener.Addr())

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
