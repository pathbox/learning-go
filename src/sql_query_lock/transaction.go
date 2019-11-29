package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 当并发超过连接池连接时，查询会完全僵死，而之前使用的连接又不得close归还到连接池中，出现像死锁一样的情况

func main() {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&timeout=3s&readTimeout=30s&writeTimeout=60s",
		"root", "", "127.0.0.1", 3306, "fintech_cms")

	db, err := sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	// db setting
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)

	transaction(db)
	c := make(chan int)
	<-c
}

func transaction(db *sql.DB) {

	sqlStr1 := `UPDATE article SET title = 'It is a start for transaction' WHERE id = 11;`
	sqlStr2 := `UPDATE article SET title = 'It is a end for transaction' WHERE id = 11;`

	tx, _ := db.Begin()
	defer tx.Rollback()
	// _, err := db.Exec(sqlStr1)
	_, err := tx.Exec(sqlStr1) // 这样就对啦！
	if err != nil {
		panic(err)
	}
	// panic("stop here======")
	_, err = tx.Exec(sqlStr2)
	if err != nil {
		panic(err)
	}
	// panic("stop here")
	fmt.Println("==done")

	tx.Commit()

}
