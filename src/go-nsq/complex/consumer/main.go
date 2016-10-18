package main

import (
	"fmt"
	"log"
	"os"

	nsq "github.com/nsqio/go-nsq"
)

var logger *log.Logger

type Handler func(*nsq.Message)

type queue struct {
	callback Handler
	*nsq.Consumer
}

func (q *queue) HandleMessage(message *nsq,Message) error {
	q.callback(message)
	return nil
}

func main() {
	fd, _ := os.OpenFile("./consumer.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	logger = log.New(fd, "", log.LstdFlags)
	config := nsq.NewConfig()
	topic := "hello world"
	channel := "test"
	c, _ := nsq.NewConsumer(topic, channel, config)
	c.SetLogger(logger, nsq.LogLevelDebug)
	q := &queue{HandleTest, c}
	c.AddHandler(q)
	addr := "127.0.0.1:4150"
	err := c.ConnectToNSQD(addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	<-q.StopChan
}

func HandleTest(msg *nsqMessage) {
	fmt.Println(string(msg.Body))
	msg.Finish()
}
