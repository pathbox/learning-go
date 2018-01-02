package main

import (
	"log"
	"time"

	dm "./dmaster"
)

func main() {
	m, err := dm.NewMaster("sd-test", []string{
		"http://127.0.0.1:2379",
	})
	if err != nil {
		log.Fatal(err)
	}
	for {
		log.Println("all ->", m.GetNodes())
		log.Println("all(strictly) ->", m.GetNodesStrictly())
		time.Sleep(time.Second * 2)
	}
}
