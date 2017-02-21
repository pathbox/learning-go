package main
import (
	"fmt"
	"github.com/scottkiss/kaca"
	"time"
)

func main() {
	producer := kaca.NewClient(":9099", "ws")
	consumer := kaca.NewClient(":9099", "ws")
	c2 := kaca.NewClient(":9099", "ws")
	c2.ConsumeMessage(func(message string) {
		fmt.Println("c2 consume =>" + message)
	})
	consumer.Sub("say")
	consumer.Sub("you")
	consumer.ConsumeMessage(func(message string) {
		fmt.Println("consume =>" + message)
	})
	time.Sleep(time.Second * time.Duration(2))
	producer.Broadcast("broadcast...")
	time.Sleep(time.Second * time.Duration(2))
}


