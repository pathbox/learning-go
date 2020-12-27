package main

import (
	"fmt"
	"sync"
	"time"
)

type Pubsub struct {
	mu     sync.RWMutex
	subs   map[string][]chan string // 每个topic 对应一个channel数组
	closed bool
}

func NewPubsub() *Pubsub {
	ps := &Pubsub{}
	ps.subs = make(map[string][]chan string)
	ps.closed = false
	return ps
}

func (ps *Pubsub) Subscribe(topic string) <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string, 1)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch // 得到该topic的 channel返回
}

func (ps *Pubsub) Publish(topic string, msg string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		return
	}

	for _, ch := range ps.subs[topic] { // 将msg 发送到该topic的所有channel
		ch <- msg
	}
}

func (ps *Pubsub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, subs := range ps.subs {
			for _, ch := range subs {
				close(ch) // 关掉所有topic的所有channel
			}
		}
	}
}

func main() {
	ps := NewPubsub()
	ch1 := ps.Subscribe("tech")
	ch2 := ps.Subscribe("travel")
	ch3 := ps.Subscribe("travel")

	listener := func(name string, ch <-chan string) {
		for i := range ch {
			fmt.Printf("[%s] got %s\n", name, i)
		}
		fmt.Printf("[%s] done\n", name)
	}

	go listener("1", ch1)
	go listener("2", ch2)
	go listener("3", ch3)

	pub := func(topic string, msg string) {
		fmt.Printf("Publishing @%s: %s\n", topic, msg)
		ps.Publish(topic, msg)
		time.Sleep(1 * time.Millisecond)
	}

	time.Sleep(50 * time.Millisecond)
	pub("tech", "tablets")
	pub("health", "vitamins")
	pub("tech", "robots")
	pub("travel", "beaches")
	pub("travel", "hiking")
	pub("tech", "drones")

	time.Sleep(50 * time.Millisecond)
	ps.Close()
	time.Sleep(50 * time.Millisecond)

}
