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

	b := req.CompareFileEtag(fileName1, filePath)

	fmt.Println("Compare File Etag:", b)

	req.HeadFile(fileName1)

	rb := req.DumpResponse(true)
	fmt.Println("check file OK:", string(rb))
	s := req.LastResponseHeader["Content-Type"]
	fmt.Println("Response Header: ", req.LastResponseHeader, s)

}
