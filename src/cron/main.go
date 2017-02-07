package main

import (
	"github.com/robfig/cron"
	"log"
)

func main() {
	i := 0
	c := cron.New()
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		i++
		log.Println("Cron runing: ", i)
	})
	c.Start()
	select {} //阻塞主线程不退出
}
