package bitcask

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type Bitcask struct {
	option   Option
	index    *index
	lock     *os.File
	oldFiles *BitFiles
	actFile  *BitFile
	mu       *sync.RWMutex
}

func New(dir string) (*Bitcask, error) {
	bf, err := newBitFile(dir)
	if err != nil {
		return nil, err
	}

	options := NewOption(dir, 0)
	bitcask := &Bitcask{
		option:   options,
		index:    newIndex(),
		lock:     lockFile,
		oldFiles: newBitFiles(),
		actFile:  bf,
		mu:       &sync.RWMutex{},
	}
	bitcask.loadIndex()
	bitcask.merge()
	return bitcask, nil
}

func (b *Bitcask) Close() {
	b.actFile.fp.Close()
	b.lock.Close()
	os.Remove(b.lock.Name())
}

func (b *Bitcask) loadIndex() {
	files, err := scanOldFiles(b.option.Die)
	if err != nil {
		panic(err)
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	for _, file := range files {
		fid, _ := getFid(file.Name())
		fp, err := os.Open(filepath.Join(b.option.Dir, file.Name()))
		if err != nil {
			continue
		}
		bitFile, err := toBitFile(fid, fp)
		if err != nil {
			continue
		}
		var offset int64 = 0
		b.oldFiles.add(fid, bitFile)
		for {
			entry, keySize, entrySize := newEntryFromBuf(fp, fid, offset)
			if entry == nil {
				break
			}

			readOffset := offset + HeaderSize
			offset += int64(entrySize)

			keyByte := make([]byte, keySize)
			if _, err := fp.ReadAt(keyByte, readOffset); err != nil {
				continue
			}
			if entry.valueSize == 0 {
				b.index.del(string(keyByte))
				continue
			}
			b.index.put(string(keyByte), entry)
		}
	}
}

func (b *Bitcask) Put(key, value []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.checkFile()

	entry, err := b.actFile.write(key, value) // write key and value 通actFile
	if err != nil {
		return err
	}
	b.index.put(string(key), entry) // update the index
	return nil
}

func (b *Bitcask) Get(key []byte) ([]byte, error) {
	entry, err := b.index.get(key)
	if err != nil {
		return nil, err
	}
	bf, err := b.checkFileState(entry.fileID)
	if err != nil {
		return nil, err
	}
	value, err := bf.read(entry.valueOffset, entry.valueSize)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (b *Bitcask) Del(key []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, err := b.index.get(key)
	if err != nil {
		return nil
	}

	b.checkFile()
	err = b.actFile.del(key)
	if err != nil {
		return err
	}
	b.index.del(string(key))
	return nil
}

func (b *Bitcask) checkFileState(fid uint32) (*BitFile, error) {
	if fid == b.actFile.fid {
		return b.actFile, nil
	}

	if bf, ok := b.oldFiles.files[fid]; ok {
		return bf, nil
	}

	return nil, errors.New("fid not exist")
}

func (b *Bitcask) checkFile() error {
	if b.actFile.offset > b.option.MaxFileSize {
		b.actFile.fp.Close()
		b.oldFiles.add(b.actFile.fid, b.actFile)

		bf, err := newBitFile(b.option.Dir) // 新生成一个file
		if err != nil {
			return err
		}

		b.actFile = bf
	}

	return nil
}
