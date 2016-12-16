package main

import (
    "fmt"
    "time"
)

func main {
    t := time.Now()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tm, _ := time.ParseInLocation("20060102150405", timestamp, loc)  // 获得准确本地时间的方法
	x := t.Sub(tm).Minutes()
    fmt.Println(t)
    fmt.Println(tm)
    fmt.Println(x)
}