package main

import (
	"github.com/scottkiss/kaca"
)

func main() {
	//use true to set check origin
	kaca.ServeWs(":9099",true)
}