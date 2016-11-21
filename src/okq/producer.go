package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/mediocregopher/okq-go/okq"
)

func main() {
	cl := okq.New("127.0.0.1:4777")
	defer cl.Close()

	for {
		err := cl.Push("super-queue", "my awesome event"+strconv.Itoa(rand.Intn(100)), okq.Normal)
		if err != nil {
			return err
		}
		fmt.Println("produce")
		time.Sleep(2 * time.Second)
	}
}
