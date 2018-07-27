package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/uaccount")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	for j := 0; j < 100; j++ {
		// go db.Ping() // 连接池几乎失效，每个goroutine都会新建一个conn连接MySQL导致大量连接产生：
		// [mysql] 2018/07/19 12:01:08 packets.go:36: read tcp 127.0.0.1:63237->127.0.0.1:3306: read: connection reset by peer
		db.Ping() // 复用了连接池连接
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Ping err: ", err)
	}
	db.SetConnMaxLifetime(-1) //If d <= 0, connections are reused forever
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	for i := 0; i < 10000; i++ {
		go func(i int) { // 下面的查询在不同的goroutine中执行，但是复用了连接池连接，在高并发时不会短时间产生大量的MySQL连接，和 go db.Ping()的情况不同,db.Ping()产生了大量的MySQL连接,that is not good
			err = db.Ping()
			if err != nil {
				fmt.Println("Ping err: ", err)
			}
		}(i)

	}
	select {}
}
