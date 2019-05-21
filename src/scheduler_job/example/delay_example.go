package main

import (
	"fmt"

	scheduler "github.com/prprprus/scheduler"
)

// Delay
func main() {
	s, err := scheduler.NewScheduler(1000)
	if err != nil {
		panic(err) // just example
	}

	jobID := s.Delay().Second(1).Do(task1, "prprprus", 23)

	// delay with 1 minute, job function without arguments
	jobID = s.Delay().Minute(1).Do(task2)

	// delay with 1 hour
	jobID = s.Delay().Hour(1).Do(task2)

	// Note: execute immediately
	jobID = s.Delay().Do(task2)

	// cancel job
	jobID = s.Delay().Day(1).Do(task2)
	err = s.CancelJob(jobID)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("cancel delay job success")
	}
}

func task1(name string, age int) {
	fmt.Printf("run task1, with arguments: %s, %d\n", name, age)
}

func task2() {
	fmt.Println("run task2, without arguments")
}
