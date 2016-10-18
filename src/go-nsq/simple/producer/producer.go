package main

import (
	"github.com/nsqio/go-nsq"
	"log"
)

func main() {
	config := nsq.NewConfig()
	w, _ := nsq.NewProducer("127.0.0.1:4150", config)

	err := w.Publish("write_test", []byte("hello world"))
	if err != nil {
		log.Panic("Could not connect")
	}

	w.Stop()
}
