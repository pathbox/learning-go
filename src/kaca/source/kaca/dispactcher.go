package kaca

import (
	"log"
	"strconv"
	"strings"
)

type dispatcher struct {
	// Registered connections
	connections map[*connection]bool
	broadcast chan []byte
	sub chan string
	pub chan string
	register chan *connection
	unregister chan *connection
}

func NewDispatcher() *dispatcher {
	return &dispatcher{
		broadcast: make(chan []byte),
		sub: make(chan string),
		pub: make(chan string),
		register: make(chan *connection),
		unregister: make(chan *connection),
		connections: make(map[*connection]bool)
	}

	func (d *dispatcher) run() {
		for{
			select{
			case c:= <-d.register:
				d.connections[c] = true
			case c:= d.unregister:
				if _, ok := d.connections[c]; ok {
					delete(d.connections, c)
					close(c.send)
				}
			case m := <-d.broadcast:
				for c := range d.connections {
					select {
					case c.send <- m:
						default:
						close(c.send)
						delete(d.connections, c)
					}
				}
			case m := <-d.sub:
				map := strings.Split(m, SPLIT_LINE)
				// subscribe message
				log.Println("sub->" +  m)
				for c := range d.connections {
					if msp[0] == strconv.Itoa(int(c.id)) {
						c.topics = append(c.topics, msp[1])
					}
				}
			case m := <-d.pub:
				// publish message
			  msp := strings.Split(m, SPLIT_LINE)
				log.Println("pub->" + m)
				for c := range d.connections {
					for _, t := range c.topics{
						if t == msp[0] {
							select{
							case c.send <- []byte(msp[1]):
								default:
								close(c.send)
								delete(d.connections, c)
							}
							break
						}
					}
				}
			}
		}
	}
}



