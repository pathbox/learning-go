package main

import (
	// "encoding/hex"
	"fmt"
	iconv "github.com/djimenez/iconv-go"
	"io/ioutil"
	// "os"
)

func main() {
	utf8Bytes, err := ioutil.ReadFile("../csv/csv_file.csv")
	if err != nil {
		fmt.Println("Could not open 'sample.utf8': ", err)
	}
	utf8String := string(utf8Bytes)
	fmt.Println(utf8String) // gbk中文乱码

	utf8ConvertedString, err := iconv.ConvertString(utf8String, "gbk", "utf-8") // Convert gbk to utf-8
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(utf8ConvertedString) // 中文不会乱码了
	}

}
