package main

import (
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	line int = 200000
	//line int = 2000
)

// 创建一个大文件
func main() {
	curDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	checkErr(err)
	fileDir := path.Join(curDir, "file-store")

	// create file
	filename := "big_input_file.txt"
	inputFile, err := os.Create(path.Join(fileDir, filename))
	checkErr(err)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := 0; i < line; i++ {
		words := r.Intn(10) // every line has at most 10 words
		var line string
		for j := 0; j < words; j++ {
			letters := 0 // every word's length is between [3, 10)
			for letters < 3 {
				letters = r.Intn(10)
			}
			var word string
			for k := 0; k < letters; k++ {
				ch := r.Intn(26)
				word += string(ch + 97)
			}
			line += word + " "
		}
		inputFile.WriteString(line + "\n")
	}
	inputFile.Close()
}
