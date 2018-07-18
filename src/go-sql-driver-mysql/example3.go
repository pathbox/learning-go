package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var (
		id   int
		name string
	)
	db, err := sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/hello") // 用户:密码@tcp(host:port)/db_name
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Ping()

	rows, err := db.Query("select id, name from users where id = ?", 1)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
