package main

import (
	"log"
	"os/exec"
)

func main() {
	// cmd := exec.Command("wkhtmltopdf", "https://qii404.me/2016/07/22/wkhtmltopdf.html", "html.pdf")

	log.Println("PDF Saving...")

	args := []string{}
	args = append(args, "--")
	args = append(args, "/usr/bin/wkhtmltopdf")
	args = append(args, "--encoding")
	args = append(args, "utf-8")
	args = append(args, "https://qii404.me/2016/07/22/wkhtmltopdf.html") // 也可以换成是html文件
	args = append(args, "htmlnew.pdf")

	// cmd := exec.Command("xvfb-run", "--", "/usr/bin/wkhtmltopdf", "--encoding", "utf-8", "https://qii404.me/2016/07/22/wkhtmltopdf.html", "htmlnew.pdf")

	cmd := exec.Command("xvfb-run", args...)

	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}

	log.Println("Done~")

}

// xvfb-run -- /usr/bin/wkhtmltopdf --encoding utf-8 https://qii404.me/2016/07/22/wkhtmltopdf.html zzz.pdf
