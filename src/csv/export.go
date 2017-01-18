package main

import (
	"encoding/csv"
	"os"
)

func main() {
	file, err := os.Create("export_csv.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	writer := csv.NewWriter(file)
	writer.Write([]string{"id", "姓名", "分数"})

	writer.Write([]string{"1", "张三", "23"})

	writer.Write([]string{"2", "李四", "24"})

	writer.Write([]string{"3", "王五", "25"})

	writer.Write([]string{"4", "赵六", "26"})

	for i := 0; i < 100000; i++ {
		writer.Write([]string{"北京市|海淀区|知春路", "多选框1", "多选框3", "状态为开启哦", "多余列", "解决中", "标准", "多余列3", "ABC,EFG", "agent@udesk.cn", "默认", "Udesk是什么1", "默认组11", "12122662223@qq.com", "你好我是简单描述", "1 1,2,3", "蓝色", "http://www.udesk.cn", "勇士", "火箭", "来个字段", "受理客服组不存在"})
	}

	writer.Flush()
}
