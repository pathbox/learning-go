package main

import (
	"log"
	"time"

	"github.com/arriqaaq/qalam"
)

func main() {
	config := qalam.NewConfig("./log.%Y%m%d.gz", time.Local, 4096, true, 10*time.Millisecond)
	c, err := qalam.NewQalam(config)
	if err != nil {
		log.Fatalln(err)
	}
	c.Writeln([]byte("pogba"))
	c.Writeln([]byte("kante"))
	c.Close()
}
