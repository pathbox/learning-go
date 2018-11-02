package main

import (
	"fmt"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	configFile = "./config.json"
)

func main() {
	filePath := "/Users/pathbox/test_file.txt"
	config, err := ufsdk.LoadConfig(configFile)
	if err != nil {
		panic(err.Error())
	}

	fileName1 := "test_file_1.txt" // 在ufile存储显示的文件名
	req, _ := ufsdk.NewFileRequest(config, nil)

	err = req.AsyncMPut(filePath, fileName1, "")
	if err != nil {
		fmt.Println("文件上传失败!!，错误信息为：", err.Error())
		//把 HTTP 详细的 HTTP response dump 出来
		fmt.Printf("%s\n", req.DumpResponse(true))
	}

	fmt.Println("Upload-Put OK:", string(req.DumpResponse(true)))

}
