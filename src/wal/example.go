package main

import (
	"github.com/tidwall/wal"
)
func main() {
	log, err := wal.Open("mylog",nil)
	// write some entries
	err = log.Write(1, []byte("first entry"))
	err = log.Write(2, []byte("second entry"))
	err = log.Write(3, []byte("third entry"))

	// read an entry
	data, err := log.Read(1)
	println(string(data))  // output: first entry
	if err != nil {
		panic(err)
	}
	f, _ := log.FirstIndex()
	l,_ := log.LastIndex()
	println(f)
	println(l)
	// close the log
	log.Close()		

}