package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func main() {
	db, err := sql.Open("mysql", "admin:password.cn@tcp(127.0.0.1:3306)/users?charset=utf8")
	if err != nil {
		fmt.Println(err)
		return
	}

	zoneID := 1
	sqlI := "INSERT INTO my_table (id, number, p, created, updated) VALUES"

	t := time.Now().Unix()
	ts := strconv.Itoa(int(t))
	zs := strconv.Itoa(zoneID)
	data := ""

	for i := 0; i < 1000; i++ { // 批量插入1000条记录
		if i == 999 { // 处理最后一条记录，末尾不需要 `,`
			bit := strconv.Itoa(i)
			item := "(" + zs + "," + bit + "," + "0," + ts + "," + ts + ")"
			data = data + item
		} else {
			bit := strconv.Itoa(i)
			item := "(" + zs + "," + bit + "," + "0," + ts + "," + ts + "),"
			data = data + item
		}
	}

	_, err = db.Exec(sqlI + data)
	if err != nil {
		fmt.Println(err)
	}
}
