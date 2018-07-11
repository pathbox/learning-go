package main

import (
	"bufio"
	"context"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/streadway/amqp"
)

var url = flag.String("url", "amqp:///", "192.168.153.41:5672")

const exchange = "test.pubsub"

type message []byte

type session struct {
	*amqp.Connection
	*amqp.Channel
}

// Close tears the connection down, taking the channel with it.
func (s session) Close() error {
	if s.Connection == nil {
		return nil
	}
	return s.Connection.Close()
}

func redial(ctx context.Context, url string) chan chan session {
	sessions := make(chan chan session) // chan 中的元素是chan

	go func() {
		sess := make(chan session)
		defer close(sessions)

		for {
			select {
			case sessions <- sess:
			case <-ctx.Done():
				log.Println("shutting down session factory")
				return
			}

			conn, err := amqp.Dial(url)
			if err != nil {
				log.Fatalf("cannot (re)dial: %v: %q", err, url)
			}

			ch, err := conn.Channel() // open one channel
			if err != nil {
				log.Fatalf("cannot create channel: %v", err)
			}
			// kind is exchange type: fanout derect topic head
			if err = ch.ExchangeDeclare(exchange, "topic", false, true, false, false, nil); err != nil {
				log.Fatalf("cannot declare fanout exchange: %v", err)
			}

			select {
			case sess <- session{conn, ch}:
			case <-ctx.Done():
				log.Println("shutting down new session")
				return
			}
		}
	}()
	return sessions
}

func publish(sessions chan chan session, messages <-chan message) {
	for session := range sessions {
		var (
			running bool
			reading = messages
			pending = make(chan message, 1)
			confirm = make(chan amqp.Confirmation, 1)
		)

		pub := <-session

		if err := pub.Confirm(false); err != nil {
			log.Printf("publisher confirms not supported")
			close(confirm) // confirms not supported, simulate by always
		} else {
			pub.NotifyPublish(confirm)
		}
		log.Printf("publishing...")

	Publish:
		for {
			var body message
			select {
			case confirmed, ok := <-confirm:
				if !ok {
					break Publish
				}
				if !confirmed.Ack {
					log.Printf("nack message %d, body: %q", confirmed.DeliveryTag, string(body))
				}
				reading = messages

			case body = <-pending:
				// routingKey := "ignored for fanout exchanges, application dependent for other exchanges"
				routingKey := "#"
				fmt.Println("publish body: ", string(body))
				err := pub.Publish(exchange, routingKey, false, false, amqp.Publishing{ // Publish 的是 exchange 和 routingKey
					Body: body,
				})
				// Retry failed delivery on the next session
				if err != nil {
					pending <- body
					pub.Close()
					break Publish
				}

			case body, running = <-reading:
				// all messages consumed
				if !running {
					return
				}
				// work on pending delivery until ack'd
				pending <- body
				reading = nil
			}
		}
	}
}

func identity() string {
	hostname, err := os.Hostname()
	h := sha1.New()
	fmt.Fprint(h, hostname)
	fmt.Fprint(h, err)
	fmt.Fprint(h, os.Getpid())
	return fmt.Sprintf("%x", h.Sum(nil))
}

func subscribe(sessions chan chan session, messages chan<- message) {
	queue := identity()
	fmt.Println("queue:", queue)
	for session := range sessions {
		sub := <-session
		if _, err := sub.QueueDeclare(queue, false, true, true, false, nil); err != nil {
			log.Printf("cannot consume from exclusive queue: %q, %v", queue, err)
			return
		}

		// routingKey := "application specific routing key for fancy toplogies"
		routingKey := "#"
		if err := sub.QueueBind(queue, routingKey, exchange, false, nil); err != nil {
			log.Printf("cannot consume without a binding to exchange: %q, %v", exchange, err)
			return
		}
		// consume from the queue
		deliveries, err := sub.Consume(queue, "#", false, true, false, false, nil) // 从 queue 中 得到publish的数据
		if err != nil {
			log.Printf("cannot consume from: %q, %v", queue, err)
			return
		}

		log.Printf("subscribed...")
		for msg := range deliveries {
			fmt.Println("receive msg: ", string(msg.Body))
			messages <- message(msg.Body)
			sub.Ack(msg.DeliveryTag, false)
		}
	}
}

func read(r io.Reader) <-chan message {
	lines := make(chan message)
	go func() {
		defer close(lines)
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			lines <- message(scan.Bytes())
		}
	}()
	return lines
}

func write(w io.Writer) chan<- message {
	lines := make(chan message)
	go func() {
		for line := range lines {
			fmt.Fprintln(w, string(line))
		}
	}()
	return lines
}

func main() {
	flag.Parse()

	ctx, done := context.WithCancel(context.Background())

	go func() {
		publish(redial(ctx, *url), read(os.Stdin))
		done()
	}()

	go func() {
		subscribe(redial(ctx, *url), write(os.Stdout))
		done()
	}()

	<-ctx.Done()
}
