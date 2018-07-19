package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/gommon/log"
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
		go func(i int) { // 下面的查询在不同的goroutine中执行，但是复用了连接池连接，在高并发时不会短时间产生大量的MySQL连接，和 db.Ping()的情况不同,db.Ping()产生了大量的MySQL连接,that is not good
			idAry := []string{"1"}
			ids := strings.Join(idAry, "','")
			sqlRaw := fmt.Sprintf(`SELECT id, resource_id, resource_type FROM t_resource WHERE resource_id IN ('%s') OR id IN ('%s')`, ids, ids)
			rows, err := db.Query(sqlRaw)

			if err != nil {
				log.Errorf("SQL t_resource error:%s", err)
			} else {
				fmt.Println("here")
				for rows.Next() {
					cols, _ := rows.Columns()
					buff := make([]interface{}, len(cols)) // 临时slice
					vals := make([]string, len(cols))      // 存数据slice
					for i, _ := range buff {
						buff[i] = &vals[i]
					}
					err = rows.Scan(buff...)
					if err != nil {
						log.Errorf("collect rows.Scan error:%s", err)
					}
					fmt.Printf("Vals:%v\n", vals)

					id := vals[0]
					resourceID := vals[1]
					resourceType := vals[2]

					fmt.Printf("id:%s, resourceID:%s, resourceType:%s--i:%d", id, resourceID, resourceType, i)
				}
			}
		}(i)

	}
	select {}
}
