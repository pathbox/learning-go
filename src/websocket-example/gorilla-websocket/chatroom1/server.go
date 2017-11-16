package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct { // 将客户端抽象为一个struct
	id     string
	socket *websocket.Conn
	send   chan []byte // 存储发送的数据
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register: // 注册上一个客户端,意味着有一个client连接连接上来了
			manager.clients[conn] = true
			jsonMessage, _ := json.Marshal(&Message{Content: "A new socket has connected."})
			manager.send(jsonMessage, conn) // 第二步
		case conn := <-manager.unregister: // 某个连接掉线,则关闭这个连接, 删除这个client conn 在register 列表
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "A socekt has disconnected"})
				manager.send(jsonMessage, conn)
			}
		case message := <-manager.broadcast: // 广播  第七步
			for conn := range manager.clients {
				select {
				case conn.send <- message: // 循环第三步
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore { // 发送的数据不要给自己发一遍
			conn.send <- message // 第三步
		}
	}
}

func (c *Client) read() {
	defer func() { // 读取结束 要关闭socket
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil { // 报错 要关闭socket
			manager.unregister <- c
			c.socket.Close()
			break
		}

		jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
		manager.broadcast <- jsonMessage // 将读取到的数据 广播到 聊天室  第六步
	}
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send: // 第四步 开始进行write操作
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			fmt.Println("websocket.TextMessage", websocket.TextMessage)
			c.socket.WriteMessage(websocket.TextMessage, message) // 第五步 客户端得到了写过来的数据,之后客户端将读到的数据又传到服务端, 然后触发 read操作
		}
	}
}

func main() {
	fmt.Println("Starting server...")

	go manager.start()
	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":12345", nil)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{ // 这一步是必要的  进行 Upgrade的操作,得到 websocket 的 conn
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(res, req, nil)

	if err != nil {
		fmt.Println(err)
		http.NotFound(res, req)
		return
	}

	client := &Client{
		id:     uuid.NewV4().String(),
		socket: conn,
		send:   make(chan []byte),
	}
	// 建立了一个socket连接,创建相应的client
	manager.register <- client // 进行注册操作, 注册操作是会阻塞的  这里是第一步
	go client.read()
	go client.write()
}
