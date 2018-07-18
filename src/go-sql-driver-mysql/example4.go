package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO foo VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // danger!
	for i := 0; i < 10; i++ {
		_, err = stmt.Exec(i)
		if err != nil {
			log.Fatal(err)
		}
	}
	stmt.Close()
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// *sql.Tx一旦释放，连接就回到连接池中，这里stmt在关闭时就无法找到连接。所以必须在Tx commit或rollback之前关闭statement
