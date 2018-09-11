package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "admin:password.cn@tcp(127.0.0.1:3306)/users?charset=utf8")
	if err != nil {
		return
	}

	sqlC := "SELECT  1 AS one FROM my_table WHERE user_id = 1 LIMIT 1;"

	r, _ := db.Query(sqlC)
	// db.QueryRow(sqlC).Scan(&isExist) // sql: no rows in result set when no record, so it is not good to do with QueryRow()

	var i bool // true or false

	for r.Next() {
		r.Scan(&i)
	}

	if i {
		fmt.Println("Result: ", i)
	} else {
		fmt.Println("False: ", i)
	}
}

// (is exists) This is a better way to check out that if the where data records exists than with count(), it is more efficient.
