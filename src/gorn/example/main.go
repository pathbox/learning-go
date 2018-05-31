package main

import (
	"fmt"
	"time"

	"github.com/roylee0704/gron"
	"github.com/roylee0704/gron/xtime"
)

type printJob struct{ Msg string }

// 定义一个struct，这个struct要实现 Run()方法，这样就实现了Job接口。Run()的具体逻辑，就是要定时执行的逻辑
func (p printJob) Run() { //
	fmt.Println(p.Msg)
}

func main() {

	var (
		daily     = gron.Every(1 * xtime.Day)
		weekly    = gron.Every(1 * xtime.Week)
		monthly   = gron.Every(30 * xtime.Day)
		yearly    = gron.Every(365 * xtime.Day)
		purgeTask = func() { fmt.Println("purge unwanted records") }
		printFoo  = printJob{"Foo"}
		printBar  = printJob{"Bar"}
	)

	c := gron.New()
	// AddFunc 底层使用Add实现的，AddFunc的第二个参数 func(){} 会通过JobFunc()适配为Job接口后，再传入Add。
	// 或者直接定义一个struct，实现Run()方法，实现Job接口，然后作为参数，直接传到Add
	c.AddFunc(gron.Every(1*time.Hour), func() {
		fmt.Println("Every 1 hour")
	})
	c.Start()

	c.AddFunc(weekly, func() { fmt.Println("Every week") })
	c.Add(daily.At("12:30"), printFoo)
	c.Start()

	// Jobs may also be added to a running Cron
	c.Add(monthly, printBar)
	c.AddFunc(yearly, purgeTask)

	// Stop the scheduler (does not stop any jobs already running).
	defer c.Stop()
}
