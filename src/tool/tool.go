package tool

import (
	"log"
	"os"
)

const (
	ConverterTmpFilePath = "/tmp/converter_tmp_file"
)

func IsBlank(s string) bool {
	if len(s) == 0 {
		return true
	}
	return false
}

func IsExist(s string) bool {
	return !IsBlank(s)
}

func MkdirTmpFile() {
	path := ConverterTmpFilePath

	if isDirExists(path) {
		log.Println(path, "-exists")
	} else {
		err := os.Mkdir(path, 0777)
		if err != nil {
			log.Println(path, "-not exists")
		} else {
			log.Println(path, "-exists")
		}
	}
}

func isDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}
