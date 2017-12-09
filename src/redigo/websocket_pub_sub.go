// Gorilla Websocket, Redigo (Redis client) and UUID
// 使用 redis的 Pub/Sub 进行服务内部的传递通讯. websocket是外部的客户端连接处理
package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var (
	gStore      *Store // gStore 是一个全局的,可以简单的认为是一个全局的存储队列
	gPubSubConn *redis.PubSubConn
	gRedisConn  = func() (redis.Conn, error) {
		return redis.Dial("tcp", ":6379")
	}
)

func init() {
	gStore = &Store{
		Users: make([]*User, 0, 1),
	}
}

type User struct {
	ID   string
	conn *websocket.Conn // 每个User使用自己的websocket连接
}

type Store struct {
	Users []*User
	sync.Mutex
}

// Message gets exchanged between users through redis pub/sub messaging
// Users may have websocket connections to different nodes and stored in
// different instances of this application
type Message struct {
	DeliveryID string `json:"id"`
	Content    string `json:"content"`
}

func (s *Store) newUser(conn *websocket.Conn) *User {
	u := &User{
		ID:   uuid.NewV4().String(),
		conn: conn,
	}

	if err := gPubSubConn.Subscribe(u.ID); err != nil {
		panic(err)
	} // 每个user 注册一个pub/sub通道, u.ID 是唯一标识

	s.Lock()
	defer s.Unlock()

	s.Users = append(s.Users, u)
	return u
}

var serverAddress = ":8080"

func main() {
	gRedisConn, err := gRedisConn() // gRedisConn全局变量就是一个func()
	if err != nil {
		panic(err)
	}
	defer gRedisConn.Close()

	gPubSubConn = &redis.PubSubConn{Conn: gRedisConn} // 开启pub/sub conn
	defer gPubSubConn.Close()

	go deliverMessages() // 创建goroutine 发送数据

	http.HandleFunc("/ws", wsHandler) // websocket handler
	log.Printf("server started at %s\n", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrader error %s\n" + err.Error())
		return
	}

	u := gStore.newUser(conn)
	log.Printf("user %s joined\n", u.ID)

	for {
		var m Message

		if err := u.conn.ReadJSON(&m); err != nil {
			log.Printf("error on ws. message %s\n", err)
		} // 将数据读取到 m

		if c, err := gRedisConn(); err != nil {
			log.Printf("error on redis conn. %s\n", err)
		} else {
			c.Do("PUBLISH", m.DeliveryID, string(m.Content)) // DeliveryID 就是对应的userID, 给这个user发送content消息
		}
	}
}

// 处理消息
func deliverMessages() {
	for { // 解析收到的数据
		switch v := gPubSubConn.Receive().(type) {
		case redis.Message: // 处理发送的消息
			log.Println("channel", v.Channel)
			gStore.findAndDeliver(v.Channel, string(v.Data)) //channel 就是注册的userID
		case redis.Subscription: // 处理注册sub 操作
			log.Printf("subscription message: %s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			log.Println("error pub/sub on connection, delivery has stopped")
			return
		}
	}
}

func (s *Store) findAndDeliver(userID string, content string) {
	m := Message{
		Content: content,
	}

	// 处理Store中的每个user // 这里其实可以用map存储
	for _, u := range s.Users {
		if u.ID == userID {
			if err := u.conn.WriteJSON(m); err != nil {
				log.Printf("error on message delivery through ws. e: %s\n", err)
			} else {
				log.Printf("Success user %s found at our store, message sent\n", userID)
			}
			return
		}
	}
	log.Printf("Fail user %s not found at our store\n", userID)
}

// > var con0 = new WebSocket('ws://localhost:8080/ws')
// > undefined
// > var con1 = new WebSocket('ws://localhost:8080/ws')
// > undefined

// > con0.onmessage = function(e) { console.log("connection 0 received message", e.data) }

// > var mes = new Object()
// > mes.id = "7c27943d-dd98–4bfe-829f-7bd9834f9f63"
// > mes.content = "hello"
// > con1.send(JSON.stringify(mes))
// > VM154:1 connection 0 received message {"id":"","content":"hello"}
