package send_big_file_with_minimal_memory

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func NormalMethod(filePath, url string) error {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	defer writer.Close()

	part, err := writer.CreateFormFile("myFile", "foo.txt")
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(part, file); err != nil {
		return err
	}

	http.Post(url, writer.FormDataContentType(), buf)
}
