package main

import (
	"log"
	"os/exec"
)

func main() {
	log.Println("WORD Saving...")

	args := []string{}
	args = append(args, "-f")
	args = append(args, "/home/user/htmlconvert/baidu.html")
	args = append(args, "-o")
	args = append(args, "/home/user/htmlconvert/baidu_html.docx")

	cmd := exec.Command("pandoc", args...)

	err := cmd.Run()

	if err != nil {
		log.Println(err)
	}

	log.Println("Done~")
}

// 使用 pandoc 进行 html => docx 的格式转换
