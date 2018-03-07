package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func main() {
	content := "<h1>hello world</h1>"
	result := htmlAll(content)
	// log.Println(result)

	createTempHtml(result)
}

func htmlAll(content string) string {
	htmlPage := `
		<html>
			<head>
				<meta http-equiv="Content-Type" content="text/html; charset=UTF-8”>
			</head>
			<body>
				%s
			</body>
		</html>`

	htmlPage = fmt.Sprintf(htmlPage, content)

	return htmlPage
}

func createTempHtml(content string) {
	fmt.Println(content)
	tempDir := "/home/user/tmp_html_file/"
	buf := &bytes.Buffer{}
	var err error
	tmp, err := ioutil.TempDir(tempDir, "temp_pdf")
	fmt.Println(tmp)
	if err != nil {
		fmt.Errorf("Error creating temp directory")
	}

	filename := fmt.Sprintf("%v/html_page.html", tmp)
	_, err = buf.WriteString(content)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(filename)
	err = ioutil.WriteFile(filename, buf.Bytes(), 0666)
	if err != nil {
		fmt.Errorf("Error writing temp file: %v", err)
	}
}
