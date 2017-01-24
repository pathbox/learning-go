package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func CreateInterrupt() {
	go func() {
		log.Println("Waiting for signal")
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
		<-c
		go func() {
			<-c
			os.Exit(1)
		}()
		log.Println("Got SIG")
		os.Exit(0)
	}()
}

func main() {
	CreateInterrupt()
	time.Sleep(10 * time.Second)
}
