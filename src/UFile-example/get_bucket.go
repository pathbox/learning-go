package main

import (
	"fmt"
	"log"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"
)

const (
	configFile = "config.json"
)

func main() {
	log.SetFlags(log.Lshortfile)
	config, err := ufsdk.LoadConfig(configFile)
	if err != nil {
		panic(err.Error())
	}
	bucketName := config.BucketName
	req, err := ufsdk.NewBucketRequest(config, nil)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(req)

	bucketList, err := req.DescribeBucket(bucketName, 0, 20, "")
	if err != nil {
		log.Println("获取 bucket 信息出错，错误信息为：", err.Error())
	} else {
		log.Println("获取 bucket list 成功，list 为", bucketList)
	}
}
