package main

import (
	"log"
	"os"
	"sync"

	"github.com/nsqio/go-nsq"
	uuid "github.com/satori/go.uuid"
)

var (
	mqkvTopic = "mqkv_topic"
)

func main() {

	wg := &sync.WaitGroup{}
	wg.Add(1)

	config := nsq.NewConfig()
	hostname, _ := os.Hostname()
	channel := hostname + ":" + uuid.NewV4().String()
	log.Printf("channel: %v\n", channel)
	// q, _ := nsq.NewConsumer(mqkvTopic, channel, config)
	q, _ := nsq.NewConsumer("write_test", "ch", config)
	q.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		log.Printf("Got Message body: %v\n", string(message.Body))
		wg.Done()
		return nil
	}))

	err := q.ConnectToNSQLookupd("10.0.3.126:4161")
	if err != nil {
		log.Panic("Could not connect")
	}

	wg.Wait()
}
