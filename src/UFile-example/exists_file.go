package main

import (
	"fmt"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	// configFile = "./config.json"
	configFile = "/Users/pathbox/code/learning-go/src/UFile-example/config.json"
)

func main() {
	config, err := ufsdk.LoadConfig(configFile)
	if err != nil {
		panic(err.Error())
	}
	req, _ := ufsdk.NewFileRequest(config, nil)
	fileName := "0001.pdf"

	err = req.HeadFile(fileName)
	if err == nil {
		fmt.Println("File is exists")
	} else {
		fmt.Println("File not exists")
		fmt.Println("File err:", err)
	}

	rb := req.DumpResponse(true)
	fmt.Println("check file OK:", string(rb))
	s := req.LastResponseHeader["Content-Length"]
	fmt.Println("Response Header: ", req.LastResponseHeader)

	fmt.Println("Etag--", s)

}
