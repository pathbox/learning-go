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
		"xxx", "xxx.cn", "xxx", 3306, "xxx")

	db, err := sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	// db setting
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)

	for i := 0; i < 100; i++ {
		fmt.Printf("Num:%d\n", i)
		go query(db, i)
	}

	c := make(chan int)
	<-c
}

func query(db *sql.DB, i int) {

	var companyID int

	sqlS := `
		SELECT company_id
		FROM t_company
		WHERE company_id= 100 limit 1
		`
	rows, err := db.Query(sqlS)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&companyID)
		sql2 := `
		SELECT company_id
		FROM t_member
		WHERE company_id= 100
		`
		fmt.Println("Start...")
		rows, err = db.Query(sql2)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			rows.Scan(&companyID)
		}
		fmt.Printf("Company ID: %d\n", companyID)
	}

	fmt.Printf("Number: %d\n", i)

}
