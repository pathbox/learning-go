package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

type Number struct {
	Content string
	AppID   uint
}

var db *gorm.DB

func init() {
	db, _ = gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/udesk_cti?charset=utf8")
	db.DB().SetMaxOpenConns(30)
	db.DB().SetMaxIdleConns(30)
	db.LogMode(true)
}

func main() {
	startHttpServer()
}

func startHttpServer() {
	http.HandleFunc("/pool", pool)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func pool(w http.ResponseWriter, r *http.Request) {
	// rows, err := db.Raw("SELECT * FROM numbers").Rows() // (*sql.Rows, error)
	// defer rows.Close()
	// checkErr(err)

	// columns, _ := rows.Columns()
	// scanArgs := make([]interface{}, len(columns))
	// values := make([]interface{}, len(columns))
	// for j := range values {
	// 	scanArgs[j] = &values[j]
	// }
	// record := make(map[string]string)
	// for rows.Next() {
	// 	//将行数据保存到record字典
	// 	err = rows.Scan(scanArgs...)
	// 	for i, col := range values {
	// 		if col != nil {
	// 			record[columns[i]] = string(col.([]byte))
	// 		}
	// 	}
	// }

	// fmt.Println(record)
	// fmt.Fprintln(w, "finish")
	number := Number{}
	// db.Where("content = ?", "18521524153").First(&number)
	// fmt.Println(number.Content, number.AppID)
	db.Model(&number).Where("content = ?", "18521524153").Update("content", "28521524153")
	db.Model(&number).Where("content = ?", "28521524153").Update("content", "18521524153")
	fmt.Println(number)
	r.Header.Set("Connection", "close")
}

// func checkErr(err error) {
// 	if err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}
// }

// ab -c 1000 -n 10000 "localhost:9090/pool"
