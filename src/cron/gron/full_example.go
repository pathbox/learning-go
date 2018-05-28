package main

import (
	"fmt"
	// "github.com/roylee0704/gron"
	// "github.com/roylee0704/gron/xtime"
)

type PrintJob struct{ Msg string }

func (p PrintJob) Run() {
	fmt.Println(p.Msg)
}

func main() {

	var (
		// schedules
		daily   = gron.Every(1 * xtime.Day)
		weekly  = gron.Every(1 * xtime.Week)
		monthly = gron.Every(30 * xtime.Day)
		yearly  = gron.Every(365 * xtime.Day)

		// contrived jobs
		purgeTask = func() { fmt.Println("purge aged records") }
		printFoo  = printJob{"Foo"}
		printBar  = printJob{"Bar"}
	)

	c := gron.New()

	c.Add(daily.At("12:30"), printFoo)
	c.AddFunc(weekly, func() { fmt.Println("Every week") })
	c.Start()

	// Jobs may also be added to a running Gron
	c.Add(monthly, printBar)
	c.AddFunc(yearly, purgeTask)

	// Stop Gron (running jobs are not halted).
	c.Stop()
}
