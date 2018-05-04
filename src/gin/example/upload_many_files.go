package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.MaxMultipartMemory = 100 << 20 // 100M
	router.Static("/", "./public")

	router.POST("/upload", UploadFiles)

	router.Run(":8080")
}

// 表单提交，上传文件
func UploadFiles(c *gin.Context) {
	name := c.PostForm("name") // 普通 参数 使用PostForm 取到
	email := c.PostForm("email")

	form, err := c.MultipartForm() // 上传文件对象 使用MultipartForm取到
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	files := form.File["files"] // 继续取出files 对象

	for _, file := range files {
		if err := c.SaveUploadedFile(file, file.Filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}
	}
	c.String(http.StatusOK, fmt.Sprintf("Uploaded successfully %d files with fields name=%s and email=%s.", len(files), name, email))

}

/*
curl -X POST http://localhost:8080/upload \
  -F "upload[]=@/Users/appleboy/test1.zip" \
  -F "upload[]=@/Users/appleboy/test2.zip" \
  -H "Content-Type: multipart/form-data"
*/
