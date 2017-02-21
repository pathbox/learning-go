package kaca


import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/url"
)

type client struct {
	id   uint64
	addr string
	path string
	conn *websocket.Conn
}

func NewClient(addr, path string) *client {
	u := url.URL{Scheme: "ws", Host: addr, Path: path}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial: ", err)
	}
	return &client{
		id: uint64(rand.Int63()),
		addr: addr,
		path: path,
		conn: c,
	}
}

func (c *client) Broadcast(message string) {
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("write:", err)
	}
}

func (c *client) Pub(topic, message string) {
	sendMsg := PUB_PREFIX + topic + SPLIT_LINE + message
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(sendMsg))
	if err != nil {
		log.Println("write:", err)
	}
}

func (c *client) Sub(topic string) {
	sendMsg := SUB_PREFIX + topic
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(sendMsg))
	if err != nil {
		log.Println("write:", err)
	}
	log.Println("sub topic :" + topic + "success")
}

func (c *client) ConsumeMessage(f func(m string)) {
	go func() {
		for {
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)
			f(string(message))
		}
	}()
}

func (c *client) Shutdown() {
	c.conn.Close()
}