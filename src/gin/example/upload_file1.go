package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	// router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Run(":8080")

}

func UploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	// Upload the file to specific dst.
	dst := "/tmp/file-name"
	c.SaveUploadedFile(file, dst)

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

/*
func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	io.Copy(out, src)
	return nil
}

curl -X POST http://localhost:8080/upload \
  -F "file=@/Users/appleboy/test.zip" \
	-H "Content-Type: multipart/form-data"

*/
