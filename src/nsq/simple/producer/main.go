package main

import (
	"log"

	"github.com/nsqio/go-nsq"
)

func main() {
	config := nsq.NewConfig()
	w, _ := nsq.NewProducer("127.0.0.1:9090", config)

	err := w.Publish("write_test", []byte("hello world test"))
	if err != nil {
		log.Panic("Could not connect")
	}

	w.Stop()
}
