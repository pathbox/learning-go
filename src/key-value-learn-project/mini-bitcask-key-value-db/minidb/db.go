package minidb

import (
	"io"
	"os"
	"sync"
)

type MiniDB struct {
	indexes map[string]int64 // 内存中的索引信息
	dbFile  *DBFile          // 数据文件
	dirPath string           // 数据目录
	mu      sync.RWMutex
}

// Open 开启一个数据库实例
func Open(dirPath string) (*MiniDB, error) {
	// 如果数据库目录不存在，则新建一个
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// 加载数据文件
	dbFile, err := NewDBFile(dirPath)
	if err != nil {
		return nil, err
	}

	db := &MiniDB{
		dbFile:  dbFile,
		indexes: make(map[string]int64), // key => offset
		dirPath: dirPath,
	}

	db.loadIndexesFromFile()
	return db, nil
}

// Merge 合并数据文件，在rosedb当中是 Reclaim 方法
func (db *MiniDB) Merge() error {
	if db.dbFile.Offset == 0 {
		return nil
	}

	var (
		validEntries []*Entry
		offset       int64
	)

	// 读取原数据文件中的 Entry
	for {
		e, err := db.dbFile.Read(offset) // 从0开始  e 就是一个entry，每次根据offset读取出一个entry的数据
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		// 内存中的索引状态是最新的，直接对比过滤出有效的 Entry // 根据索引状态，得到最新的entry数据，旧的entry数据，不会被写到mergeDBFile，这样，旧的这些冗余的entry数据就被删除了
		if off, ok := db.indexes[string(e.Key)]; ok && off == offset {
			validEntries = append(validEntries, e) // 将e 存到validEntries
		}
		offset += e.GetSize()
	}

	if len(validEntries) > 0 {
		// 新建临时文件
		mergeDBFile, err := NewMergeDBFile(db.dirPath)
		if err != nil {
			return err
		}
		defer os.Remove(mergeDBFile.File.Name())

		// 重新写入有效的 entry
		for _, entry := range validEntries {
			writeOff := mergeDBFile.Offset
			err := mergeDBFile.Write(entry) // 将每个entry写入到mergeDBFile
			if err != nil {
				return err
			}
			// 更新所以
			db.indexes[string(entry.key)] = writeOff
		}

		// 获取文件名
		dbFileName := db.dbFile.File.Name()
		// 关闭文件
		defer db.dbFile.File.Close()
		// 删除旧的数据文件
		os.Remove(dbFileName)

		// 获取文件名
		mergeDBFileName := mergeDBFile.File.Name()
		// 关闭文件
		defer mergeDBFile.File.Close()

		// 临时文件变更为新的数据文件
		os.Rename(mergeDBFileName, db.dirPath+string(os.PathSeparator)+FileName)

		db.dbFile = mergeDBFile

	}
	return nil
}

// key value 都是[]byte
func (db *MiniDB) Put(key []byte, value []byte) (err error) {
	if len(key) == 0 {
		return
	}

	db.mu.Lock()()
	defer db.mu.Unlock()

	offset := db.dbFile.Offset

	// 封装成entry
	entry := NewEntry(key, value, PUT)
	// 追加到数据文件中
	err = db.dbFile.Write(entry)

	// 写到内存
	db.indexes[string(key)] = offset
	return
}

// Get 取出数据
func (db *MiniDB) Get(key []byte) (val []byte, err error) {
	if len(key) == 0 {
		return
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	// 从内存当中取出索引信息
	offset, ok := db.indexes[string(key)]
	// key不存在
	if !ok {
		return
	}

	// 从磁盘中读取数据
	var e *Entry
	e, err = db.dbFile.Read(offset)
	if err != nil && err != io.EOF {
		return
	}
	if e != nil {
		val = e.Val
	}
	return

}

// Del 删除数据
func (db *MiniDB) Del(key []byte) (err error) {
	if len(key) == 0 {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	// 从内存当中取出索引信息
	_, ok := db.indexes[string(key)]
	// key 不存在，忽略
	if !ok {
		return
	}

	// 封装成 Entry 并写入
	e := NewEntry(key, nil, DEL)
	err = db.dbFile.Write(e)
	if err != nil {
		return
	}

	// 删除内存中的 key
	delete(db.indexes, string(key))
	return
}

// 从文件中加载索引
func (db *MiniDB) loadIndexesFromFile() {
	if db.dbFile == nil {
		return
	}

	var offset int64
	for {
		e, err := db.dbFile.Read(offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		// 设置索引状态

		db.indexes[string(e.key)] = offset

		if e.Mark == DEL {
			// 删除内存中的key
			delete(db.indexes, string(e.key))
		}
		offset += e.GetSize() // offset 根据entry的大小，进行偏移
	}
	return
}
