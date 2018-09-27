package main

import (
	"github.com/arriqaaq/qalam"
	"log"
	"time"
)

func main() {
	config := qalam.NewConfig("./log.%Y%m%d.gz", time.Local, 1, true, 10*time.Millisecond)
	c, err := qalam.NewQalam(config)
	if err != nil {
		log.Fatalln(err)
	}
	c.Writeln([]byte("pogba"))
	c.Writeln([]byte("kante"))
	c.Close()
}