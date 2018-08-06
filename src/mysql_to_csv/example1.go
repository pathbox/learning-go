package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

var orgID = flag.String("org_id", "", "org_id value for query")
var optTime = flag.Int("opt_time", 0, "opt_time value for query")

type OpLog struct {
	UserEmail      string `db:"user_email"`
	Api            string `db:"api"`
	ObjectInfoList string `db:"object_info_list"`
	OptTime        int64  `db:"opt_time"`
	RemoteIP       string `db:"remote_ip"`
}

func main() {
	flag.Parse()

	if len(*orgID) == 0 || *optTime == 0 {
		fmt.Println("org_id, opt_time args wrong")
		return
	}

	fmt.Printf("org_id: %s, opt_time: %d\n", *orgID, *optTime)

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
	timeLayout := "2006-01-02 15:04:05"
	sqlRaw := fmt.Sprintf("SELECT user_email, api, object_info_list, opt_time, remote_ip FROM t_user_opt_log WHERE org_id = '%s' AND opt_time >= %d limit 1000;", *orgID, *optTime)
	// for {
	err = db.Select(&logs, sqlRaw)
	if err != nil {
		fmt.Println("db select error: ", err)
		return
	}

	for _, log := range logs {
		dateStr := time.Unix(log.OptTime, 0).Format(timeLayout)
		line := []string{"测试项目名", log.UserEmail, log.Api, log.ObjectInfoList, dateStr, log.RemoteIP}
		writer.Write(line)
	}
	writer.Flush()
	// }
}

// ./example1 -org_id=xxx -opt_time=xxx
// ./example1 -org_id=orgnaization_1 -opt_time=1533031711
// 如果将 writer.Flush() 放到for循环中，CPU 到100% 而内存只占0.2%，说明内存只用了一次flush所需要的内存，不会导致爆内存，对内存来说是安全的
