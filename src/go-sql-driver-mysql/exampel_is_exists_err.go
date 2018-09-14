package main

import (
	"database/sql"
	"errors"
	"fmt"
)

func main() {
	r, err := GetBitLevel(100, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(r)
}

// 如果记录不存在就返回，则可以使用这种方式，不存在会报错然后返回，如果记录不存在，不要立即返回，而是继续下面的逻辑，那这种方法是不合适的
func GetBitLevel(b, z int, a interface{}) (string, error) {
	// 分zoneID 等于0和不等于0
	if zoneID == 0 {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/user")
		if err != nil {
			return "", err
		}
		ag := a.(int)
		var e int
		sqlG := `SELECT 1 AS one FROM t
						 WHERE b = ? AND z = 0 AND a = 0
						 LIMIT 1;`
		sqlA := `SELECT 1 AS one FROM t
						 WHERE b = ? AND z = 0 AND a  = ?
						 LIMIT 1;`
		if err = db.QueryRow(sqlG, b).Scan(&e); err != nil {
			return "", err
		}
		if e != 0 {
			return "G", nil
		}
		if err = db.QueryRow(sqlA, b, ag).Scan(&e); err != nil {
			return "", err
		}
		if e != 0 {
			return "A", nil
		}

		return "", errors.New(" WRONG")
	}
}
