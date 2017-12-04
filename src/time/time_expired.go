package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func main() {
	timestampStr := "1512388480"
	expireStr := "300"
	expireInt, _ := strconv.Atoi(expireStr)
	timestampInt, _ := strconv.Atoi(timestampStr)
	timestampNow := int(time.Now().Unix())
	val := timestampNow - timestampInt
	fmt.Println("val", val)
	fmt.Println(timestampInt, timestampNow, expireInt)
	if val > expireInt {
		log.Panic("Login token is expired")
	} else {
		log.Println("Go")
	}
}
