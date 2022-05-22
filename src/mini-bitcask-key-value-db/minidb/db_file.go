package minidb

import "os"

const FileName = "minidb.data"
const MergeFileName = "minidb.data.merge"

// DBFile 数据文件定义
type DBFile struct {
	File   *os.File
	Offset int64
}

func newInternal(fileName string) (*DBFile, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(fileName)
	if err != nil {
		return nil, err
	}
	return &DBFile{Offset: stat.Size(), File: file}, nil // stat.Size() 是初始文件的偏移，这之后才是数据
}

// NewDBFile 创建一个新的数据文件
func NewDBFile(path string) (*DBFile, error) {
	fileName := path + string(os.PathSeparator) + FileName
	return newInternal(fileName)
}

// NewMergeDBFile 新建一个合并时的数据文件
func NewMergeDBFile(path string) (*DBFile, error) {
	fileName := path + string(os.PathSeparator) + MergeFileName
	return newInternal(fileName)
}

// Read 从 offset 处开始读取
func (df *DBFile) Read(offset int64) (e *Entry, err error) {
	buf := make([]byte, entryHeaderSize)
	if _, err = df.File.ReadAt(buf, offset); err != nil { // 1.读取header的值
		return
	}
	if e, err = Decode(buf); err != nil {
		return
	}

	offset += entryHeaderSize
	if e.KeySize > 0 {
		key := make([]byte, e.KeySize)
		if _, err = df.File.ReadAt(key, offset); err != nil { // 2.读取key的值
			return
		}
		e.Key = key
	}

	offset += int64(e.KeySize)
	if e.ValueSize > 0 {
		value := make([]byte, e.ValueSize)
		if _, err = df.File.ReadAt(value, offset); err != nil { // 读取value的值
			return
		}
		e.Value = value
	}
	return
	// 可以看到 offset 是每次读取数据的起始点
}

// Write 写入 Entry
func (df *DBFile) Write(e *Entry) (err error) {
	enc, err := e.Encode() // 将entry encode为[]byte数据
	if err != nil {
		return err
	}
	_, err = df.File.WriteAt(enc, df.Offset) // 从df.Offset起始点开始写入enc数据
	df.Offset += e.GetSize
	return
}
