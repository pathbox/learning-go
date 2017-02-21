package kaca

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
	SUB_PREFIX = "__sub:"
	PUB_PREFIX = "__pub:"
	maxTopics = 100
	SPLIT_LINE = "_:_"
)

var disp = NewDispatcher()

type connection struct {
	// websocket connection.
	ws *websocket.Conn
	send chan []byte
	topics []string
	cid uint64
}

func (c *connection) deliver() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select{
		case message, ok := <-c.send:
			if !ok {
				c.sendMsg(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.sendMsg(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.sendMsg(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *connection) dispatch() {
	defer func() {
		disp.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil})
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		msg := string(message)
		if strings.Contains(msg, SUB_PREFIX){
			topic := strings.Split(msg, SUB_PREFIX)[1]
			disp.sub <- strconv.Itoa(int(c.cid)) + SPLIT_LINE + topic
		} else if strings.Contains(msg, PUB_PREFIX) {
			topic_msg := strings.Split(msg, PUB_PREFIX)[1]
			disp.pub <- topic_msg
		} else {
			disp.broadcast <- message
		}
	}
}

func (c *connection) dispatch() {
	defer func(){
		disp.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetReadLimit(maxmessageSize)
	c.ws.SetPongHandler(func(string) error {c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil})
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		msg := string(message)
		if strings.Contains(msg, SUB_PREFIX) {
			topic := strings.Split(msg, SUB_PREFIX)[1]
			disp.sub <- strconv.Itoa(int(c.cid)) + SPLIT_LINE + topic
		} else if strings.Contains(msg, PUB_PREFIX) {
			topic_msg := strings.Split(msg, PUB_PREFIX)[1]
			disp.pub <- topic_msg
		} else {
			disp.broadcast <- message
		}
	}
}

func (c *connection) sendMsg(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{cid: uint64(rand.Int63()), send: make(chan []byte, 256), ws: ws, topic: make([]string, maxTopics)}
	disp.register <- c
	go c.dispatch()
	c.deliver()
}

func serveWsCheckOrigin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{cid: uint64(rand.Int63()), send: make(chan []byte, 256), ws: ws, topics: make([]string, maxTopics)}
	disp.register <- c
	go c.dispatch()
	c.deliver()
}


















