package main

import (
	"fmt"
	"log"
	"os"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// fasthttp 对文件上传的部分没有做大修改，使用和 net/http 一样
func httpHandler(ctx *fasthttp.RequestCtx) {
	// 这里直接获取到 multipart.FileHeader, 需要手动打开文件句柄
	f, err := ctx.FormFile("file")
	if err != nil {
		ctx.SetStatusCode(500)
		fmt.Println("get upload file error: ", err)
		return
	}
	fh, err := f.Open()
	if err != nil {
		fmt.Println("open upload file error:", err)
		ctx.SetStatusCode(500)
		return
	}
	defer fh.Close()

	// 打开文件句柄
	fp, err := os.OpenFile("saveto.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("open saving file error:", err)
		ctx.SetStatusCode(500)
		return
	}
	defer fp.Close()
	// 把ｆｈ　保存到　ｆｐ
	if _, err = io.Copy(fp, fh); err != nil {
		fmt.Println("save upload file error:", err)
		ctx.SetStatusCode(500)
		return
	}
	ctx.Write([]byte("save file successfully!"))
}

func main() {
	//　使用　fasthttprouter 创建路由
	router := fasthttprouter.New()
	router.GET("/", httpHandler)
	log.Println(fasthttp.ListenAndServe(":9090", router.Handler))
}
