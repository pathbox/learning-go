package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	dw "./dworker"
)

func main() {
	name := flag.String("name", fmt.Sprintf("%d", time.Now().Unix()), "des")
	extInfo := "lhq-demo..."

	flag.Parse()
	w, err := dw.NewWorker("sd-test", *name, extInfo, []string{
		"http://127.0.0.1:2379",
	})
	if err != nil {
		log.Fatal(err)
	}
	w.Register()
	log.Println("name ->", *name, "extInfo ->", extInfo)

	go func() {
		time.Sleep(time.Second * 20)
		w.Unregister()
	}()

	for {
		log.Println("isActive ->", w.IsActive())
		log.Println("isStop ->", w.IsStop())
		time.Sleep(time.Second * 2)
		//服务退出
		if w.IsStop() {
			return
		}
	}
}
