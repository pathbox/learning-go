package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	// Channel name to use with redis
	Channel = "chat"
)

var (
	waitingMessage, availableMessage []byte
	waitSleep                        = time.Second * 10
)

func init() {
	var err error
	waitingMessage, err = json.Marshal(message{
		Handle: "system",
		Text:   "Waiting for redis to be available. Messaging won't work until redis is available",
	})
	if err != nil {
		panic(err)
	}
	availableMessage, err = json.Marshal(message{
		Handle: "system",
		Text:   "Redis is now available & messaging is now possible",
	})
	if err != nil {
		panic(err)
	}
}

// redisReceiver receives messages from Redis and broadcasts them to all
// registered websocket connections that are Registered.
type redisReceiver struct {
	pool *redis.Pool

	messages       chan []byte
	newConnections chan *websocket.Conn
	rmConnections  chan *websocket.Conn
}

// newRedisReceiver creates a redisReceiver that will use the provided
// rredis.Pool.
func newRedisReceiver(pool *redis.Pool) redisReceiver {
	return redisReceiver{
		pool:           pool,
		messages:       make(chan []byte, 1000), // 1000 is arbitrary
		newConnections: make(chan *websocket.Conn),
		rmConnections:  make(chan *websocket.Conn),
	}
}

func (rr *redisReceiver) wait(_ time.Time) error {
	rr.broadcast(waitingMessage)
	time.Sleep(waitSleep)
	return nil
}

// run receives pubsub messages from Redis after establishing a connection.
// When a valid message is received it is broadcast to all connected websockets
func (rr *redisReceiver) run() error {
	l := log.WithField("channel", Channel)
	conn := rr.pool.Get()
	defer conn.Close()
	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe(Channel)
	go rr.connHandler()
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			l.WithField("message", string(v.Data)).Info("Redis Message Received")
			if _, err := validateMessage(v.Data); err != nil {
				l.WithField("err", err).Error("Error unmarshalling message from Redis")
				continue
			}
			rr.broadcast(v.Data)
		case redis.Subscription:
			l.WithFields(logrus.Fields{
				"kind":  v.Kind,
				"count": v.Count,
			}).Println("Redis Subscription Received")
		case error:
			return errors.Wrap(v, "Error while subscribed to Redis channel")
		default:
			l.WithField("v", v).Info("Unknown Redis receive during subscription")
		}
	}
}

// broadcast the provided message to all connected websocket connections.
// If an error occurs while writting a message to a websocket connection it is
// closed and deregistered.
func (rr *redisReceiver) broadcast(msg []byte) {
	rr.messages <- msg
}

// register the websocket connection with the receiver.
func (rr *redisReceiver) register(conn *websocket.Conn) {
	rr.newConnections <- conn
}

// deRegister the connection by closing it and removing it from our list.
func (rr *redisReceiver) deRegister(conn *websocket.Conn) {
	rr.rmConnections <- conn
}

func (rr *redisReceiver) connHandler() {
	conns := make([]*websocket.Conn, 0)
	for {
		select {
		case msg := <-rr.messages:
			for _, conn := range conns {
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.WithFields(logrus.Fields{
						"data": msg,
						"err":  err,
						"conn": conn,
					}).Error("Error writting data to connection! Closing and removing Connection")
					conns = removeConn(conns, conn)
				}
			}
		case conn := <-rr.newConnections:
			conns = append(conns, conn)
		case conn := <-rr.rmConnections:
			conns = removeConn(conns, conn)
		}
	}
}

func removeConn(conns []*websocket.Conn, remove *websocket.Conn) []*websocket.Conn {
	var i int
	var found bool
	for i = 0; i < len(conns); i++ {
		if conns[i] == remove {
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("conns: %#v\nconn: %#v\n", conns, remove)
		panic("Conn not found")
	}
	copy(conns[i:], conns[i+1:]) // shift down
	conns[len(conns)-1] = nil    // nil last element
	return conns[:len(conns)-1]  // truncate slice
}

// redisWriter publishes messages to the Redis CHANNEL
type redisWriter struct {
	pool     *redis.Pool
	messages chan []byte
}

func newRedisWriter(pool *redis.Pool) redisWriter {
	return redisWriter{
		pool:     pool,
		messages: make(chan []byte, 10000),
	}
}

// run the main redisWriter loop that publishes incoming messages to Redis.
func (rw *redisWriter) run() error {
	conn := rw.pool.Get()
	defer conn.Close()

	for data := range rw.messages {
		if err := writeToRedis(conn, data); err != nil {
			rw.publish(data) // attempt to redeliver later
			return err
		}
	}
	return nil
}

func writeToRedis(conn redis.Conn, data []byte) error {
	if err := conn.Send("PUBLISH", Channel, data); err != nil {
		return errors.Wrap(err, "Unable to publish message to Redis")
	}
	if err := conn.Flush(); err != nil {
		return errors.Wrap(err, "Unable to flush published message to Redis")
	}
	return nil
}

// publish to Redis via channel.
func (rw *redisWriter) publish(data []byte) {
	rw.messages <- data
}
