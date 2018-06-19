package main

import (
	"bytes"
	"io"
	"os"
	"time"
)

func main() {
	proverbs := new(bytes.Buffer) // &bytes.Buffer{}
	proverbs.WriteString("Channels orchestrate mutexes serialize\n")
	proverbs.WriteString("Cgo is not Go\n")
	proverbs.WriteString("Errors are values\n")

	pipeReader, pipeWriter := io.Pipe()

	go func() {
		defer pipeWriter.Close()
		time.Sleep(3 * time.Second)
		io.Copy(pipeWriter, proverbs)
	}()

	io.Copy(os.Stdout, pipeReader) // 阻塞一直等到 pipe中有数据进来
	pipeReader.Close()

}
