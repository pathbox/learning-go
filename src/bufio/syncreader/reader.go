package syncreader

import (
	"bufio"
	"io"
	"os"
)

func ReadFile(filePath string, bufSize int) (chan string, error) {
	outChan := make(chan string, bufSize)
	file, err := os.Open(filePath)
	if err != nil {
		return make(chan string), err
	}

	fileReader := bufio.NewScanner(file)

	go func() {
		defer file.Close()
		for fileReader.Scan() {
			outChan <- fileReader.Text()
		}
		close(outChan)
	}()
	return outChan, nil
}

func ReadString(r io.Reader, bufSize int) (chan string, error) {
	outChan := make(chan string, bufSize)
	bufReader := bufio.NewScanner(r)

	go func() {
		for bufReader.Scan() {
			outChan <- bufReader.Text()
		}
		close(outChan)
	}()
	return outChan, nil
}

func ReadByte(r io.Reader, bufSize int) (chan []byte, error) {
	outChan := make(chan []byte, bufSize)
	bufReader := bufio.NewScanner(r)

	go func() {
		for bufReader.Scan() {
			outChan <- bufReader.Bytes()
		}
		close(outChan)
	}()
	return outChan, nil
}

// 单独起了一个goroutine 去源 io.Reader中读取数据到缓冲outChan，主goroutine再从outChan中读取数据。spawns a goroutine to read the file, and sends the lines over the returned channel
