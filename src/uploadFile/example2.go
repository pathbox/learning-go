package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

// 表单中增加enctype="multipart/form-data"
// 服务端调用r.ParseMultipartForm,把上传的文件存储在内存和临时文件中
// 使用r.FormFile获取文件句柄，然后对文件进行存储等处理。
func main() {

	target_url := "http://localhost:9090/upload"
	filename := "/Users/pathbox/code/chonghuafei.png"
	postFile(filename, target_url)
}

func postFile(filename, target_url string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// 关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// 打开文件句柄
	fh, err := os.Open(filename) // 打开要上传的文件
	if err != nil {
		fmt.Println("error opening file or file is not exits")
		return err
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}
	contentType := bodyWriter.FormDataContentType()
	fmt.Println(contentType)
	bodyWriter.Close()
	resp, err := http.Post(target_url, contentType, bodyBuf)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}
