package main

import (
	"log"

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

	log.Println("公有空间文件下载 URL 是：", req.GetPublicURL(fileName))
	log.Println("私有空间文件下载 URL 是：", req.GetPrivateURL(fileName, 1)) // 我的bucket_name 是私有空间，所以公有空间的URL是用不了的
	// 私有空间文件的下载地址是可用的，不传过期时间的话，默认应该是30分钟,过期时间参数好像没效果?有效，据说有时间差延迟是因为缓存
	fs, _ := req.PrefixFileList("test", "", 0)
	log.Println("File list: ", fs)
}
