package main

import (
	"fmt"
)

type chanWriter struct {
	ch chan byte
}

func newChanWriter() *chanWriter {
	return &chanWriter{make(chan byte, 1024)}
}

func (w *chanWriter) Chan() <-chan byte {
	return w.ch
}

func (w *chanWriter) Write(p []byte) (int, error) {
	n := 0
	for _, b := range p {
		w.ch <- b
		n++
	}
	return n, nil
}

func (w *chanWriter) Close() error {
	close(w.ch)
	return nil
}

func main() {
	writer := newChanWriter()
	go func() {
		// 如果把writer.Close()注释了，报 Stream me!fatal error: all goroutines are asleep - deadlock!
		defer writer.Close()

		writer.Write([]byte("Stream "))
		writer.Write([]byte("me!"))
	}()

	for c := range writer.Chan() { // 会阻塞等待  Write goroutine,直到 writer.Close(), 则range退出循环

		fmt.Printf("%c", c)
	}
	fmt.Println()
	fmt.Println("done")
}
