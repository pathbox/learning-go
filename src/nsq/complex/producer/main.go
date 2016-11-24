package main

import (
	"io/ioutil"
	"log"

	nsq "github.com/nsqio/go-nsq"
)

var (
	logger *log.Logger
)

func main() {
	logger = log.New(ioutil.Discard, "", log.LstdFlags)
	config := nsq.NewConfig()
	addr := "127.0.0.1:9090"
	w, _ := nsq.NewProducer(addr, config)
	defer w.Stop()
	w.SetLogger(logger, nsq.LogLevelDebug)

	msgCount := 10
	var testData [][]byte
	for i := 0; i < msgCount; i++ {
		testData = append(testData, []byte("multipublish_test_case"))
	}

	topicName := "test"
	responseChan := make(chan *nsq.ProducerTransaction)
	err := w.MultiPublishAsync(topicName, testData, responseChan, "test", 1)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	trans := <-responseChan
	if trans.Error != nil {
		logger.Fatalf(trans.Error.Error())
	}
	if trans.Args[0].(string) != "test" {
		logger.Fatalf(`proxied arg "%s" != "test"`, trans.Args[0].(string))
	}
	if trans.Args[1].(int) != 1 {
		logger.Fatalf(`proxied arg %d != 1`, trans.Args[1].(int))
	}
}
