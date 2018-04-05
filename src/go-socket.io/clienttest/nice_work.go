package main

import (
	"log"
	"runtime"
	"strconv"
	// "time"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type Channel struct {
	Channel string `json:"channel"`
}

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// func sendJoin(c *gosocketio.Client) {
// 	log.Println("Acking /join")
// 	result, err := c.Ack("/join", Channel{"main"}, time.Second*5)
// 	if err != nil {
// 		log.Fatal(err)
// 	} else {
// 		log.Println("Ack result to /join: ", result)
// 	}
// }

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer recover()
 
	for i := 0; i < 3000; i++ {
		// time.Sleep(20 * time.Microsecond)
		c, err := gosocketio.Dial(
			gosocketio.GetUrl("127.0.0.1", 7002, false),
			transport.GetDefaultWebsocketTransport())
		if err != nil {
			log.Println(err)
		}

		args := make(map[string]string)
		args["uid"] = strconv.Itoa(i)
		args["client_id"] = args["uid"]
		args["castle"] = ""
		c.Emit("register", args)

		// err = c.On("/message", func(h *gosocketio.Channel, args Message) {
		// 	log.Println("--- Got chat message: ", args)
		// })
		if err != nil {
			log.Println(err)
		}

		err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
			log.Println("Disconnected")
		})
		if err != nil {
			log.Println(err)
		}

		err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
			log.Println("Connected", c.Id(), "--", i)
		})
		if err != nil {
			log.Println(err)
		}

		// time.Sleep(1 * time.Second)

		// go sendJoin(c)
		// go sendJoin(c)
		// go sendJoin(c)
		// go sendJoin(c)
		// go sendJoin(c)

		// time.Sleep(600 * time.Second)
		// defer c.Close()
	}
	select {}

	// log.Println(" [x] Complete")
}
