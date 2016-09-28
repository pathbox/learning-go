package main

import(
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  // _ "github.com/go-sql-driver/mysql"
  "fmt"
)


// import _ "github.com/jinzhu/gorm/dialects/mysql"
// import _ "github.com/jinzhu/gorm/dialects/postgres"
// import _ "github.com/jinzhu/gorm/dialects/sqlite"
// import _ "github.com/jinzhu/gorm/dialects/mssql"

func main() {
  db, err := gorm.Open("mysql", "root:@/say_morning_development?charset=utf8&parseTime=True&loc=Local")
  if err != nil{
    panic(err)
  }
  fmt.Println(db)
  defer db.Close()
}
