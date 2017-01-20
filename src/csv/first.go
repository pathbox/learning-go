package main

import (
	"encoding/csv"
	"fmt"
	iconv "github.com/djimenez/iconv-go"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("csv_file.csv") // 打开想要读取的文件
	if err != nil {
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)                      // 创建reader 读取器
	converter, _ := iconv.NewConverter("gbk", "utf-8") // 将读取的数据由gbk编码转为UTF-8编码
	for {                                              // 循环读取，直到结尾
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		joinString := strings.Join(record, " ")
		output, _ := converter.ConvertString(joinString)
		fmt.Println(output)
	}
}
