package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/download", zipHandler)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func zipHandler(w http.ResponseWriter, r *http.Request) {
	zipName := "ZipTest.zip"
	// 设置rw的header信息中的ctontent-type，对于zip可选以下两种
	// w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Type", "application/zip")
	// 设置w的header信息中的Content-Disposition为attachment类型
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", zipName))

	// 向w中写入zip文件
	err := getZip(w)
	if err != nil {
		log.Fatal(err)
	}
}

func getZip(w io.Writer) error {
	// NewWriter 中的 w 是zip压缩生成后的output writer，将http.ResponseWriter 作为这个output writer，加上 Header做相关信息的配置，让浏览器能识别到下载zip文件，这样就成功了
	zipW := zip.NewWriter(w)
	defer zipW.Close()

	for i := 0; i < 5; i++ {
		// 向zip中添加文件
		f, err := zipW.Create(strconv.Itoa(i) + ".txt")
		if err != nil {
			return err
		}
		// 向文件中写入文件内容
		_, err = f.Write([]byte(fmt.Sprintf("Hello file %d", i)))
		if err != nil {
			return err
		}
	}
	return nil
}
