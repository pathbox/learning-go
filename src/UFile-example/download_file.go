package main

import (
	"fmt"
	"os"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	configFile1 = "./config.json"
)

func main() {
	config, err := ufsdk.LoadConfig(configFile1)
	if err != nil {
		panic(err.Error())
	}

	fileName := "test_file_1.txt" // 在ufile存储显示的文件名
	req, _ := ufsdk.NewFileRequest(config, nil)

	f, err := os.Create("/Users/pathbox/download_pdf_file.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = req.DownloadFile(f, fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	s, _ := f.Stat()
	fmt.Println("File Size: ", s.Size())
}
