package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	content := []byte("temporary file's content")
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir) // 删除临时目录以及临时目录下的所有文件

	tmpfn := filepath.Join(dir, "tmpfile")
	if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
		log.Fatal(err)
	}

}
