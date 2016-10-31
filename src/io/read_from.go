package main

import (
	"bufio"
	"os"
)

func main() {
	file, err := os.Open("../http/http_client.go")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	writer := bufio.NewWriter(os.Stdout)
	writer.ReadFrom(file)
	writer.Flush()
}
