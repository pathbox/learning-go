package main

import (
	// "database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	url := "root:@tcp(127.0.0.1:3306)/mars?charseutf8"
	db, err := sqlx.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	sql := "SELECT title,content FROM posts where match(content) against('没有' IN NATURAL LANGUAGE MODE);"

	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}

	var title, content string
	for rows.Next() {
		err = rows.Scan(&title, &content)
		if err != nil {
			panic(err)
		}
		fmt.Println("Title: ", title)
		fmt.Println("Content Len: ", len(content))
	}

}
