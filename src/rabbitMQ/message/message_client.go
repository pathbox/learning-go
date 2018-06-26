package message

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type IMessagingClient interface {
	ConnectToBroker(connectionString string)
	Publish(msg []byte, exchangeName string, exchangeType string) error
	PublishOnQueue(msg []byte, queueName string) errpr
	Subscribe(exchangeName string, exchangeType string, consumerName string, handlerFunc func(amqp.Delivery)) error
	SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery)) error
	Close()
}

type MessagingClient struct {
	conn *amqp.Connection
}

// connect amqp broker
func (m *MessagingClient) ConnectToBroker(connectionString string) {
	if connectionString == "" {
		panic("Cannot initialize connection to broker, connectionString not set. Have you initialized?")
	}

	var err error
	m.conn, err = amqp.Dial(fmt.Sprintf("%s/", connectionString))
	if err != nil {
		panic("Failed to connect to AMQP compatible broker at: " + connectionString)
	}
}

// message => exchange => queue
func (m *MessagingClient) Publish(body []byte, exchangeName string, exchangeType string) error {
	if m.conn == nil {
		panic("Tried to send message before connection was initialized. Don't do that.")
	}
	ch, err := m.conn.Channel() // 从这个conn中得到一个channel,一个conn可以有多个channel
	defer ch.Close()
	//1、channel的声明配置
	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // derect topic header or fanout
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	failOnError(err, "Failed to register an Exchange")
	// 2、Declare a queue that will be created if not exists with some args
	queue, err := ch.QueueDeclare(
		"",    // our queue name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	//3、bind channel and queue
	err = ch.QueueBind(
		queue.Name,   // name of the queue
		exchangeName, // bindingKey, now routing key is exchangeName
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	)

	// 4.Publishes a message onto the queue
	err = ch.Publish(
		exchangeName, // exchange
		exchangeName, // routing key      q.Name
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Body: body, // our JSON body as []byte
		})
	// type Publishing struct {
	// 	Headers         Table
	// 	ContentType     string
	// 	ContentEncoding string
	// 	DeliveryMode    uint8
	// 	Priority        uint8
	// 	CorrelationId   string
	// 	ReplyTo         string
	// 	Expiration      string
	// 	MessageId       string
	// 	Timestamp       time.Time
	// 	Type            string
	// 	UserId          string
	// 	AppId           string
	// 	Body            []byte
	// }
	fmt.Printf("A message was sent: %v", body)
	return err
}

// without QueueBind step, message => queue 消息直接发送到queue，不经过exchange
func (m *MessagingClient) PublishOnQueue(body []byte, queueName string) error {
	if m.conn == nil {
		panic("Tried to send message before connection was initialized. Don't do that.")
	}
	ch, err := m.conn.Channel() // Get a channel from the connection
	defer ch.Close()

	queue, err := ch.QueueDeclare( // Declare a queue that will be created if not exists with some args
		queueName, // our queue name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	// Publishes a message onto the queue.
	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body, // our json body as []byte
		})
	fmt.Printf("A message was sent to queue %v: %v", queueName, body)
	return err
}

// 通过 exchange subscribe queue
func (m *MessagingClient) Subscribe(exchangeName, exchangeType, consumerName string, handlerFunc func(amqp.Delivery)) error {

	ch, err := m.conn.Channel()
	failOnError(err, "Failed to open a channel")
	// defer ch.Close()
	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false, // delete when complete
		false, // internal
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to register an Exchange")

	log.Printf("declared Exchange, declaring Queue (%s)", "")
	queue, err := ch.QueueDeclare(
		"",    // name of the queue, 这是自定义的，只要是唯一的就可以
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to register an Queue")

	log.Printf("declared Queue (%d messages, %d consumers), binding to Exchange (key '%s')",
		queue.Messages, queue.Consumers, exchangeName)

	err = ch.QueueBind(
		queue.Name,   // name of the queue
		exchangeName, // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	msgs, err := ch.Consume(
		queue.Name,   // queue
		consumerName, // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)

	failOnError(err, "Failed to register a consumer")
	// type Delivery struct {
	// 	Acknowledger    Acknowledger
	// 	Headers         Table
	// 	ContentType     string
	// 	ContentEncoding string
	// 	DeliveryMode    uint8
	// 	Priority        uint8
	// 	CorrelationId   string
	// 	ReplyTo         string
	// 	Expiration      string
	// 	MessageId       string
	// 	Timestamp       time.Time
	// 	Type            string
	// 	UserId          string
	// 	AppId           string
	// 	ConsumerTag     string
	// 	MessageCount    uint32
	// 	DeliveryTag     uint64
	// 	Redelivered     bool
	// 	Exchange        string
	// 	RoutingKey      string
	// 	Body            []byte
	// }
	// 上面根据规则，订阅玩exchange queue之后，新开一个goroutine，监听 <-chan amqp.Delivery， 得到Delivery 传入自定义的handlerFunc中，一般操作Body数据
	go consumeLoop(msgs, handlerFunc)
	return nil
}

// 直接 subscribe queue
func (m *MessagingClient) SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery)) error {
	ch, err := m.conn.Channel()
	failOnError(err, "Failed to open a channel")

	log.Printf("Declaring Queue (%s)", queueName)
	queue, err := ch.QueueDeclare(
		queueName, // name of the queue
		false,     // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	failOnError(err, "Failed to register an Queue")
	// 得到 <-chan amqp.Deliver
	msgs, err := ch.Consume(
		queue.Name,   // queue
		consumerName, // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)

	failOnError(err, "Failed to register a consumer")

	go consumeLoop(msgs, handlerFunc)
	return nil
}

func (m *MessagingClient) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
}

// 上面根据规则，订阅玩exchange queue之后，新开一个goroutine，监听 <-chan amqp.Delivery， 得到Delivery 传入自定义的handlerFunc中，一般操作Body数据
func consumeLoop(deliveries <-chan amqp.Delivery, handleFunc func(d amqp.Delivery)) {
	for d := range deliveries {
		// Invoke the handlerFunc func we passed as parameter.
		handleFunc(d)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
