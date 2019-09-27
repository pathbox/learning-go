package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
			log.Fatalf("%s: %s", msg, err)
			panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Fialed to connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	// 申明一个队列
	queue, err := channel.QueueDeclare(
		"task_queue",
		true, // durable 持久性
		false, // delete when unused 不自动删除该队列
		false, // exclusive 如果是true，连接一断开 队列就会删除
		false, // no-wait
		nil, // arguments
	)
	failOnError(err, "Fialed to declare a queue")
	body := bodyFrom(os.Args)
	// 发布
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType: "text/plain", // application/json
		Body: []byte(body),
	}
	err = channel.Publish(
		"topic", //exchange 默认模式 
		queue.Name, // routing key 默认路由到同名队列：task_name
		false, 
		false, 
		msg,
	)
	failOnError(err, "Falied to publish a message")
	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
					s = "hello"
	} else {
					s = strings.Join(args[1:], " ")
	}
	return s
}