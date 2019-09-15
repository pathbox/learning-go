package main

import (
        "fmt"
        "log"
        "strconv"

        "github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
        if err != nil {
                log.Fatalf("%s: %s", msg, err)
                panic(fmt.Sprintf("%s: %s", msg, err))
        }
}

func fib(n int) int {
        if n == 0 {
                return 0
        } else if n == 1 {
                return 1
        } else {
                return fib(n-1) + fib(n-2)
        }
}
// rpc模式是两个队列之间互相应答
func main() {
	// 拨号
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 声明匿名队列
	q, err := ch.QueueDeclare(
					"rpc_queue", // name
					false,       // durable
					false,       // delete when usused
					false,       // exclusive
					false,       // no-wait
					nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 公平分发 没有这个则round-robbin：https://segmentfault.com/a/1190000004492447
	err = ch.Qos(
					1,     // prefetch count
					0,     // prefetch size
					false, // global
	)
	failOnError(err, "Failed to set QoS")

	// 消费，等待请求
	msgs, err := ch.Consume(
					q.Name, // queue
					"",     // consumer
					false,  // auto-ack
					false,  // exclusive
					false,  // no-local
					false,  // no-wait
					nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		//请求来了
		for d := range msgs {
						n, err := strconv.Atoi(string(d.Body))
						failOnError(err, "Failed to convert body to integer")

						log.Printf(" [.] fib(%d)", n)
						
						// 计算
						response := fib(n)

						// 回答
						err = ch.Publish(
										"",        // exchange
										d.ReplyTo, // routing key  回答队列
										false,     // mandatory
										false,     // immediate
										amqp.Publishing{
											ContentType:   "text/plain",
											CorrelationId: d.CorrelationId,  序列号
											Body:          []byte(strconv.Itoa(response)),
										})
						failOnError(err, "Failed to publish a message")


						// 确认回答完毕
						d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<forever
}