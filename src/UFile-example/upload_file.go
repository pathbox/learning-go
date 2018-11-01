package main

import (
	"fmt"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	configFile1 = "./config.json"
)

func main() {
	filePath := "./test_file.txt"
	config, err := ufsdk.LoadConfig(configFile1)
	if err != nil {
		panic(err.Error())
	}

	fileName1 := "test_file_1.txt" // 在ufile存储显示的文件名
	req, _ := ufsdk.NewFileRequest(config, nil)

	err = req.PutFile(filePath, fileName1, "")
	if err != nil {
		fmt.Println("文件上传失败!!，错误信息为：", err.Error())
		//把 HTTP 详细的 HTTP response dump 出来
		fmt.Printf("%s\n", req.DumpResponse(true))
	}

	fmt.Println("Upload-Put OK:", string(req.DumpResponse(true)))

	fileName2 := "test_file_2.txt"
	err = req.PostFile(filePath, fileName2, "")
	if err != nil {
		fmt.Println("文件上传失败!!，错误信息为：", err.Error())
		//把 HTTP 详细的 HTTP response dump 出来
		fmt.Printf("%s\n", req.DumpResponse(true))
	}

	fmt.Println("Upload-Put OK:", string(req.DumpResponse(true)))

}
