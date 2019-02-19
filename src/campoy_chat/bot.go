package main

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/campoy/chat/markov"
)

var chain = markov.NewChain(2)

type bot struct {
	io.ReadCloser
	out io.Writer
}

func newBot() io.ReadWriteCloser {
	r, out := io.Pipe()
	return bot{r, out}
}

func (b bot) Write(p []byte) (int, error) {
	log.Printf("bot received: %s", p)
	if len(bytes.TrimSpace(p)) > 0 {
		go b.speak()
	}
	return len(p), nil
}

func (b bot) speak() {
	time.Sleep(time.Second)
	msg := chain.Generate(10)
	_, err := b.out.Write([]byte(msg))
	if err != nil {
		log.Printf("could not write from bot: %v", err)
	}
}