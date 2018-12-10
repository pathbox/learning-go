package send_big_file_with_minimal_memory

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func NiceMethod(filePath, url string) error {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer m.Close()

		part, err := m.CreateFormFile("myFile", "foo.txt")
		if err != nil {
			retuen err
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}

		defer file.Close()

		if _, err = io.Copy(part, file); err != nil {
			return err
		}
	}()

	http.Post(url, m.FormDataContentType(), r)
}

/*
If you dump the request above, the header reads:
POST / HTTP/1.1
...
Transfer-Encoding: chunked
Accept-Encoding: gzip
Content-Type: multipart/form-data; boundary=....
User-Agent: Go-http-client/1.1
*/
