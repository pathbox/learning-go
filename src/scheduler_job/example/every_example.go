package main

import (
	"fmt"

	scheduler "github.com/prprprus/scheduler"
)

func main() {
	s, err := scheduler.NewScheduler(1000)
	if err != nil {
		panic(err) // just example
	}

	// Specifies time to execute periodically
	jobID = s.Every().Second(45).Minute(20).Hour(13).Day(23).Weekday(3).Month(6).Do(task1, "prprprus", 23)
	jobID = s.Every().Second(15).Minute(40).Hour(16).Weekday(4).Do(task2)
	jobID = s.Every().Second(1).Do(task1, "prprprus", 23)

	// Note: execute immediately
	jobID = s.Every().Do(task2)

	jobID = s.Every().Second(1).Minute(1).Hour(1).Do(task2)
	err = s.CancelJob(jobID)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("cancel periodically job success")
	}
}
