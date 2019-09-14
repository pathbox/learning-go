package main

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
  failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 指定连接到队列:
	queue, err := ch.QueueDeclare(
		"task_name",
		true, // durable 和生产者设置保持一致
		false, // delete when unused
		false, // 如果是true，连接一断开 队列就会删除
		false,
		nil,
	)
	failOnError(err,"")

	err = ch.Qos(
		1, // prefetch count
		0, 
		false, 
	)
	failOnError(err,"")

	// 消费者根据路由的队列名
	msgs, err := ch.Consume(
		queue.Name, 
		"cpmsumer name", 
		false, // auto-ask
		false, //exclusive
		false, // no-local
		false, // no-wait
		nil, //args 
	)
	failOnError(err, "")
	forever := make(chan bool)

	go func() { // 异步从msgs 消费者中消费msg
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			// 确认消息被收到！！如果为真的，那么同在一个channel,在该消息之前未确认的消息都会确认，适合批量处理
      d.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever // 不让consumer main goroutine结束
}