package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func (c *PostGetController) sendPostFormFile(url string, filename string) error {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)
	// boundary默认会提供一组随机数，也可以自己设置。
	body_writer.SetBoundary("Pp7Ye2EeWaFDdAY")
	//boundary :=  body_writer.Boundary()

	// 1. 要上传的数据
	body_writer.WriteField("key1", "value-string1")
	body_writer.WriteField("key2", fmt.Sprintf("%d", 45))
	// 2. 内存中的文件1，FormFile1
	_, err := body_writer.CreateFormFile("filekey1", "filename.txt")
	if err != nil {
		fmt.Printf("创建FormFile1文件信息异常！", err)
		return err
	}
	f1_content := "内存中文件1的内容"
	body_buf.Write([]byte(f1_content))
	// 3. 读取文件
	_, errr := body_writer.CreateFormFile("filekey2", filename)
	if errr != nil {
		fmt.Printf("创建FormFile2文件信息异常！", err)
		return errr
	}
	fb2, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("打开文件异常!", err)
		return err
	}
	body_buf.Write(fb2)
	// 结束整个消息body
	body_writer.Close()

	//
	req_reader := io.MultiReader(body_buf)
	req, err := http.NewRequest("POST", url, req_reader)
	if err != nil {
		fmt.Printf("站点相机上传图片，创建上次请求异常！异常信息:", err)
		return err
	}
	// 添加Post头
	req.Header.Set("Connection", "close")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Content-Type", body_writer.FormDataContentType())
	req.ContentLength = int64(body_buf.Len())
	fmt.Printf("发送消息长度:", req.ContentLength)
	// 发送消息
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取回应消息异常:", err)
	}
	fmt.Printf("发送回应数据:", string(body))
	return nil
}

// 接收上传的文件
func (this *PostGetController) DoRecePostFormFiel() {
	// 读取参数信息
	value1 := this.GetString("key1")
	value2, _ := this.GetInt64("key2")
	_, fh1, ferr := this.GetFile("filekey1")
	_, fh2, ferr := this.GetFile("filekey2")

	//
	fmt.Printf("key1: ", value1)
	fmt.Printf("key2: ", value2)

	// 保存文件1到目录
	fmt.Printf("filekey1-Filename: ", fh1.Filename)
	ferr = this.SaveToFile("filekey1", fmt.Sprintf("./%s", fh1.Filename))
	if ferr != nil {
		fmt.Printf("保存文件1失败:", ferr.Error())
	}

	// 保存文件2到目录,文件名称前面加file2只是修改保存文件的名称。
	fmt.Printf("filekey2-Filename: ", fh2.Filename)
	ferr = this.SaveToFile("filekey2", fmt.Sprintf("./file2_%s", fh2.Filename))
	if ferr != nil {
		fmt.Printf("保存文件2失败:", ferr.Error())
	}
}
