package main

import (
	kafka "gopkg.in/Shopify/sarama.v2"
)

func main() {
	consumer, err := kafka.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("my topic",0, kafka.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
    if err := partitionConsumer.Close(); err != nil {
        log.Fatalln(err)
			}
	}()

// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
	ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}
	log.Printf("Consumed: %d\n", consumed)
}
