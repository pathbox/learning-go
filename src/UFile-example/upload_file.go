package main

import (
	"fmt"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	configFile1 = "config.json"
)

func main() {
	filePath := ""
	config, err := ufsdk.LoadConfig(configFile1)
	if err != nil {
		panic(err.Error())
	}
	req, _ := ufsdk.NewFileRequest(config, nil)
	err = req.PutFile(filePath, "keyName", "")
	if err != nil {
		fmt.Println("文件上传失败!!，错误信息为：", err.Error())
		//把 HTTP 详细的 HTTP response dump 出来
		fmt.Printf("===%s\n", req.DumpResponse(true))
	}
}
