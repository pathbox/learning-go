package gorm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// DB contains information for current db connection

type DB struct {
	Value             interface{}
	Error             error
	RowsAffected      int64
	callbacks         *Callback
	db                sqlCommon
	parent            *DB
	search            *search
	logMode           int
	logger            logger
	dialect           Dialect
	singularTable     bool
	source            string
	values            map[string]interface{}
	joinTableHandlers map[string]JoinTableHandler
	blockGlobalUpdate bool
}

// Open initialize a new db connection, need to import driver first, e.g:
//
//     import _ "github.com/go-sql-driver/mysql"
//     func main() {
//       db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
//     }
// GORM has wrapped some drivers, for easier to remember driver's import path, so you could import the mysql driver with
//    import _ "github.com/jinzhu/gorm/dialects/mysql"
//    // import _ "github.com/jinzhu/gorm/dialects/postgres"
//    // import _ "github.com/jinzhu/gorm/dialects/sqlite"
//    // import _ "github.com/jinzhu/gorm/dialects/mssql"
