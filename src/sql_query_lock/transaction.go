package main

import (
	"database/sql"
	"fmt"
	"strings"
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
}

func transaction(db *sql.DB) {
	sqlQ := `SELECT id FROM article limit 10`
	sqlStr1 := `UPDATE article SET title = 'It is a start for transaction' WHERE id = 11;`
	sqlStr2 := `UPDATE article SET title = 'It is a end for transaction' WHERE id = 11;`
	rows, _ := db.Query(sqlQ)
	defer rows.Close()
	idAry := make([]string, 0)
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
		}
		idAry = append(idAry, id)
		fmt.Println("===", id)
	}

	fmt.Println("===", idAry)

	// sqlQ2 := `SELECT title FROM article WHERE id IN(?)` // This is not right
	sqlQ2 := fmt.Sprintf("SELECT title FROM article WHERE id IN(%s)", strings.Join(idAry, ",")) // This is right
	fmt.Println("sqlQ2:", sqlQ2)
	rows2, err := db.Query(sqlQ2)
	if err != nil {
		fmt.Println(rows2.Err())
	}
	defer rows2.Close()

	for rows2.Next() {
		var title string
		err = rows2.Scan(&title)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("XXX", title)
	}

	return

	tx, _ := db.Begin()
	defer tx.Rollback()
	// _, err := db.Exec(sqlStr1)
	_, err = tx.Exec(sqlStr1) // 这样就对啦！
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
