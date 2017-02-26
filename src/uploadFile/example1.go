package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/upload", upload)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func upload(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile") // 解析表单提交的文件，存入内存，file为所上传文件的指针
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintln(w, "%v", handler.Header)
		fmt.Println("====================")
		fmt.Fprintln(w, "%v", handler.Filename)
		f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666) // 创建或打开了一个路径的文件句柄
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file) // 将内存中的文件数据保存到目的地
		w.Write([]byte("Upload Success!"))
	}
}

// 文件handler是multipart.FileHeader,里面存储了如下结构信息

//     type FileHeader struct {
//         Filename string
//         Header   textproto.MIMEHeader
//         // contains filtered or unexported fields
//     }
