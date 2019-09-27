package main

import (
        "fmt"
        "log"

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
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 同样要申明交换机
	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// 新建队列，如果这个队列没名字，会随机生成一个名字
	q, err := ch.QueueDeclare(
					"logs_queue",    // name
					false, // durable
					false, // delete when usused
					true,  // exclusive  表示连接一断开，这个队列自动删除
					false, // no-wait
					nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 队列和交换机绑定，即是队列订阅了发到这个交换机的消息
	err = ch.QueueBind(
					q.Name, // queue name  队列的名字
					"",     // routing key  广播模式不需要这个
					"logs", // exchange  交换机名字
					false,
					nil)
	failOnError(err, "Failed to bind a queue")


	// 开始消费消息，可开多个订阅方，因为队列是临时生成的，所有每个订阅方都能收到同样的消息
	msgs, err := ch.Consume(
					q.Name, // queue  队列名字
					"",     // consumer
					true,   // auto-ack  自动确认
					false,  // exclusive
					false,  // no-local
					false,  // no-wait
					nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
					for d := range msgs {
									log.Printf(" [x] %s", d.Body)
					}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}