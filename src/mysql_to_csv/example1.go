package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

type OpLog struct {
	UserEmail      string `db:"user_email"`
	Api            string `db:"api"`
	ObjectInfoList string `db:"object_info_list"`
	OptTime        int    `db:"opt_time"`
	RemoteIP       string `db:"remote_ip"`
}

func main() {
	file, err := os.Create("/Users/pathbox/code/learning-go/src/mysql_to_csv/export_csv.csv") // 准备好导出文件
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(file)                                             // 初始化 csv writer
	err = writer.Write([]string{"项目名", "操作人", "API", "操作资源", "操作时间", "操作IP"}) // header 数据
	if err != nil {
		panic(err)
	}
	url := "root:@tcp(127.0.0.1:3306)/ulog?charseutf8"
	db, err := sqlx.Open("mysql", url)
	if err != nil {
		panic(err)
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Ping error: ", err)
		return
	}

	logs := []OpLog{}
	sqlRaw := fmt.Sprintf("SELECT user_email, api, object_info_list, opt_time, remote_ip FROM t_user_opt_log WHERE org_id = '%s' AND opt_time >= %d;", "orgnaization_1", 1533031711)
	err = db.Select(&logs, sqlRaw)
	if err != nil {
		fmt.Println("db select error: ", err)
		return
	}
	fmt.Println(len(logs))
	for _, log := range logs {
		opTime := strconv.Itoa(log.OptTime)
		line := []string{"测试项目名", log.UserEmail, log.Api, log.ObjectInfoList, opTime, log.RemoteIP}
		writer.Write(line)
	}
	writer.Flush()
}
