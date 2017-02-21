// pub/sub client

package main

import (
	"fmt"
	"time"
	"github.com/scottkiss/kaca"
)

func main() {
	producer := kaca.NewClient(":9099", "ws")
	consumer := kaca.NewClient(":9099", "ws")
	consumer.Sub("say")
	consumer.Sub("you")
	consumer.ConsumeMessage(func(message string) {
		fmt.Println("consumer => " + message)
	})
	time.Sleep(time.Second * time.Duration(2))
	producer.Pub("you", "world")
	producer.Pub("say", "hello")
	time.Sleep(time.Second * time.Duration(2))
}

