package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

func main() {
	// redisServerAddr, err := serverAddr()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	redisServerAddr := "0.0.0.0:6379"

	ctx, cancel := context.WithCancel(context.Background())

	err := listenPubSubChannels(ctx, redisServerAddr,
		func() error {
			go publish()
			return nil
		},

		func(channel string, message []byte) error {
			fmt.Printf("channel: %s, message: %s\n", channel, message)

			// For the purpose of this example, cancel the listener's context
			// after receiving last message sent by publish().
			if string(message) == "goodbye" {
				cancel()
			}
			return nil
		}, "c1", "c2")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func listenPubSubChannels(ctx context.Context, redisServerAddr string,
	onStart func() error,
	onMessage func(channel string, data []byte) error,
	channels ...string) error {
	const healthCheckPeriod = time.Minute

	c, err := redis.Dial("tcp", redisServerAddr,
		redis.DialReadTimeout(healthCheckPeriod+10*time.Second),
		redis.DialWriteTimeout(10*time.Second))

	if err != nil {
		return err
	}

	defer c.Close()

	psc := redis.PubSubConn{Conn: c} // pub sub conn

	if err := psc.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		return err
	} // sub  these channels,可以批量subscribe 多个 channel

	done := make(chan error, 1)

	// Start a goroutine to receive notifications from the server.
	go func() {
		for {
			switch n := psc.Receive().(type) {
			case error:
				done <- n
				return
			case redis.Message:
				if err := onMessage(n.Channel, n.Data); err != nil {
					done <- err
					return
				}
			case redis.Subscription:
				switch n.Count {
				case len(channels):
					// Notify application when all channels are subscribed.
					if err := onStart(); err != nil {
						done <- err
						return
					}
				case 0:
					// Return from the goroutine when all channels are unsubscribed.
					done <- nil
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(healthCheckPeriod)
	defer ticker.Stop()

loop:
	for err == nil {
		select {
		case <-ticker.C:
			// Send ping to test health of connection and server. If
			// corresponding pong is not received, then receive on the
			// connection will timeout and the receive goroutine will exit.
			if err = psc.Ping(""); err != nil {
				break loop
			}
		case <-ctx.Done():
			break loop
		case err := <-done:
			return err
		}
	}

	// 有错误，Signal the receiving goroutine to exit by unsubscribing from all channels.

	psc.Unsubscribe()

	return <-done

}

func publish() {
	c, err := dial()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	c.Do("PUBLISH", "c1", "hello")
	c.Do("PUBLISH", "c2", "world")
	c.Do("PUBLISH", "c1", "goodbye")
}

func dial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, nil
	}
	return c, nil
}
