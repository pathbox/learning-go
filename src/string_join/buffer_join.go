package main

import (
	"bytes"
	"fmt"
)

func main() {
	var buffer bytes.Buffer
	s := []string{"appID", "@", "queueID"}
	for _, item := range s {
		buffer.WriteString(item)
	}
	fmt.Println(buffer.String())

	r := stringJoin("1", "2", "3")
	fmt.Println(r)
}

func stringJoin(args ...string) string {
	var buffer bytes.Buffer
	for _, s := range args {
		buffer.WriteString(s)
	}
	return buffer.String()
}
