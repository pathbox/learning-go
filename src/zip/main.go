package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	CompressZip()   // 压缩
	DeCompressZip() // 解压缩
}

func CompressZip() {
	const dir = "../img/"
	// 获取文件列表
	f, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}

	fzip, _ := os.Create("img.zip")
	w := zip.NewWriter(fzip)
	defer w.Close()

	for _, file := range f {
		fw, _ := w.Create(file.Name()) // 在压缩包中创建每个文件
		fileContent, err := ioutil.ReadFile(dir + file.Name())
		if err != nil {
			fmt.Println(err)
		}
		n, err := fw.Write(fileContent) // 将压缩包外的文件内容写入到压缩包内
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(n)
	}
}

func DeCompressZip() {
	const File = "img.zip"
	const dir = "img/"

	os.Mkdir(dir, 0777)
	cf, err := zip.OpenReader(File) // 读取zip文件
	if err != nil {
		fmt.Println(err)
	}

	defer cf.Close()
	for _, file := range cf.File {
		rc, err := file.Open()
		if err != nil {
			fmt.Println(err)
		}

		f, err := os.Create(dir + file.Name) // 创建一个文件
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		n, err := io.Copy(f, rc) // 将压缩包中的每个文件 复制到新建的文件,从而达到解压缩的目的
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(n)
	}

}
