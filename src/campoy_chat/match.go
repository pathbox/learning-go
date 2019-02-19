package main

import (
	"fmt"
	"io"
	"log"
	"time"
)

var partners = make(chan io.ReadWriteCloser)

func match(conn io.ReadWriteCloser) {
	fmt.Fprintln(conn, "Looking for a partner ...")
	select {
	case partners <- conn:
		// the other goroutine won and we can finish
	case p := <-partners:
		chat(conn, p)
	case <-time.After(5 * time.Second):
		log.Println("timeout, chat to bot")
		chat(newBot(), conn)
	}
}

func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "We found a partner")
	fmt.Fprintln(b, "We found a partner")

	errc := make(chan error, 1)
	go copy(a, chain.SpyOn(b), errc)
	go copy(b, chain.SpyOn(a), errc)
	if err := <-errc; err != nil {
		log.Printf("Error chatting: %v", err)
	}
	a.Close()
	b.Close()
}

func copy(a io.Writer, b io.Reader, errc chan<- error) {
	_, err := io.Copy(a, b)
	errc <- err
}
