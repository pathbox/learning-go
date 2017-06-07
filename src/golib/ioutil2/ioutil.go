package ioutil2

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

func WriteFileAtomic(filename string, data []byte, perm os.FileMode) error {
	dir, name := path.Dir(filename), path.Base(filename)
	f, err := ioutil.TempFile(dir, name) // 在定义的目录，产生一个临时文件
	if err != nil {
		return err
	}
	n, err := f.Write(data) // 写入数据
	f.Close()
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	} else {
		err = os.Chmod(f.Name(), perm)
	}
	if err != nil {
		os.Remove(f.Name()) //删除临时文件
		return err
	}
	return os.Rename(f.Name(), filename) // 将临时文件，修改为定义的文件名

}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
