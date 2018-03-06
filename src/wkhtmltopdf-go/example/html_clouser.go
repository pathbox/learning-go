package main

import (
	"fmt"
	"log"
)

func main() {
	content := "<h1>hello world</h1>"
	result := htmlAll(content)
	log.Println(result)
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
