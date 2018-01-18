package main

import (
	"time"

	workers "github.com/Scalingo/go-workers"
)

func main() {
	workers.Configure(map[string]string{
		"process": "client1",
		"server":  "localhost:6379",
	})
	workers.EnqueueIn("myqueue", "MyRubyWorker", 2*60, []string{"in two minutes"})

	now := time.Now()
	hoursTo9 := (time.Duration(9 - now.Hour())) * time.Hour
	at := now.AddDate(0, 0, 1).Truncate(time.Hour).Add(hoursTo9).Unix()

	workers.EnqueueAt("myqueue", "MyGoWorker", at, []string{"tomorrow at 9"})
}
