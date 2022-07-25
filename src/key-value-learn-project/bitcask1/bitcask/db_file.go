package bitcask

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type BitFile struct {
	fp     *os.File
	fid    uint32
	offset uint64
}

func newBitFile(dir string) (*BitFile, error) {
	bf := &BitFile{}
	fp, err := bf.openFile(dir)
	if err != nil {
		return nil, err
	}
	bf.fp = fp

	return bf, nil
}

func (bf *BitFile) populateFilesMap(dir string) (uint32, error) {
	files, err := scanOldFiles(dir)
	if err != nil {
		return 0, err
	}

	found := make(map[uint32]struct{})
	var maxFid uint32 = 0
	for _, file := range files {
		fid, err := getFid(file.Name())
		if err != nil {
			return 0, err
		}
		if _, ok := found[fid]; ok {
			return 0, errors.New("Duplicate file found.")
		}
		found[fid] = struct{}{}
		if maxFid < fid {
			maxFid = fid
		}
	}
	return maxFid, nil
}

func (bf *BitFile) write(key, value []byte) (*entry, error) {
	ts := uint32(time.Now().Unix())

	keySize := uint32(len(key))
	valueSize := uint32(len(value))
	entrySize := getSize(keySize, valueSize)
	buf, _ := encode(key, value, keySize, valueSize, ts, entrySize)

	offset := bf.offset + uint64(HeaderSize+keySize)

	_, err := bf.fp.WriteAt(buf, int64(bf.offset))
	if err != nil {
		return nil, err
	}

	bf.offset += uint64(entrySize)

	entry := newEntry(bf.fid, valueSize, offset, uint64(ts))
	return entry, nil
}

func (bf *BitFile) read(offset uint64, size uint32) ([]byte, error) {
	return read(bf.fp, int64(offset), size)
}

func (bf *BitFile) del(key []byte) error {
	ts := uint32(time.Now().Unix())
	keySize := uint32(len(key))
	var valueSize uint32 = 0
	entrySize := getSize(keySize, valueSize)
	buf, _ := encode(key, nil, keySize, valueSize, ts, entrySize)

	_, err := bf.fp.WriteAt(buf, int64(bf.offset))
	if err != nil {
		return err
	}

	bf.offset += uint64(entrySize)

	return nil
}

func (bf *BitFile) openFile(dir string) (*os.File, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	file, err := bf.newFile(dir)
	if err != nil {
		return nil, err
	}
	fp, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return fp, nil
}

func (bf *BitFile) newFile(dir string) (string, error) {
	lastFid, err := bf.populateFilesMap(dir)
	if err != nil {
		return "", err
	}

	// create a new file
	bf.fid = lastFid + 1
	fid := bf.newFid()

	dataFilePath := newFilePath(dir, fid)
	return dataFilePath, nil
}

func (bf *BitFile) newFid() string {
	return fmt.Sprintf("%06d", bf.fid)
}

func read(fp *os.File, offset int64, size uint32) ([]byte, error) {
	buf := make([]byte, size)
	if _, err := fp.ReadAt(buf, offset); err != nil {
		return nil, err
	}
	return buf, nil
}

func toBitFile(fid uint32, fp *os.File) (*BitFile, error) {
	stat, err := fp.Stat()
	if err != nil {
		return nil, err
	}

	bf := &BitFile{
		fp:     fp,
		fid:    fid,
		offset: uint64(stat.Size()),
	}

	return bf, nil
}

func getFid(name string) (uint32, error) {
	fsz := len(name)
	fid, err := strconv.ParseUint(name[:fsz-5], 10, 32)
	if err != nil {
		return 0, errors.New("Unable to parse file id.")
	}

	return uint32(fid), nil
}

func scanOldFiles(dir string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.New("Unable to open dir.")
	}
	var entry []os.DirEntry
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ext) {
			continue
		}
		entry = append(entry, file)
	}
	return entry, err
}

const ext = ".data"

func newFilePath(dir, fid string) string {
	return fmt.Sprintf("%s%s%s%s", dir, string(os.PathSeparator), fid, ext)
}

type BitFiles struct {
	files map[uint32]*BitFile
	mu    *sync.RWMutex
}

func newBitFiles() *BitFiles {
	return &BitFiles{
		files: make(map[uint32]*BitFile),
		mu:    &sync.RWMutex{},
	}
}

func (bf *BitFiles) add(fid uint32, fp *BitFile) {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	bf.files[fid] = fp
}

const (
	lockFileName = "bitcask.lock"
	mergeFileExt = ".merge"
)

func lock(dir string) (*os.File, error) {
	return os.OpenFile(filepath.Join(dir, lockFileName), os.O_EXCL|os.O_CREATE|os.O_RDWR, os.ModePerm)
}

func newEntryFromBuf(fp *os.File, fid uint32, offset int64) (*entry, uint32, uint32) {
	buf, err := read(fp, offset, HeaderSize)
	if err != nil {
		if err == io.EOF {
			return nil, 0, 0
		}
	}
	ts := binary.BigEndian.Uint32(buf[4:8])
	keySize := binary.BigEndian.Uint32(buf[8:12])
	valueSize := binary.BigEndian.Uint32(buf[12:HeaderSize])

	entrySize := getSize(keySize, valueSize)

	entry := newEntry(fid, valueSize, uint64(offset)+uint64(HeaderSize+keySize), uint64(ts))
	return entry, keySize, entrySize
}

func newMergeFileName(dir string, fid uint32) string {
	return fmt.Sprintf("%s%s%d%s", dir, string(os.PathSeparator), fid, mergeFileExt)
}

func newMergeFile(mf string) (*os.File, error) {
	fp, err := os.OpenFile(mf, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return fp, nil
}

func (b *Bitcask) merge() {
	log.Println("Start merge old files.")
	files, err := scanOldFiles(b.option.Dir)
	if err != nil {
		return
	}
	for _, file := range files {
		log.Println("File name:", file.Name())
		if filepath.Base(b.actFile.fp.Name()) == file.Name() { //skip active file
			continue
		}

		oldFilePath := filepath.Join(b.option.Dir, file.Name())
		fp, err := os.Open(oldFilePath)
		if err != nil {
			continue
		}

		info, err := fp.Stat()
		if err != nil {
			log.Println("Err check file size:", err)
			continue
		}
		if info.Size() == 0 {
			log.Printf("Removing old empy file:%s\n", file.Name())
			os.Remove(oldFilePath)
			continue
		}

		fid, _ := getFid(file.Name())
		mergeFileName := newMergeFileName(b.option.Dir, fid)
		mergeFp, err := newMergeFile(mergeFileName)

		var (
			offset      int64 = 0
			mergeOffset int64 = 0
		)
		for {
			entry, keySize, entrySize := newEntryFromBuf(fp, fid, offset)
			if entry == nil {
				break
			}

			readOffset := offset + HeaderSize
			offset += int64(entrySize)
			keyByte, err := read(fp, readOffset, keySize)
			if err != nil {
				continue
			}

			//check if the key was deleted
			e, err := b.index.get(keyByte)
			if e == nil || entry.valueSize == 0 {
				b.index.del(string(keyByte))
				continue
			}
			valByte, err := read(fp, readOffset+int64(keySize), entry.valueSize)
			if err != nil {
				continue
			}

			buf, _ := encode(keyByte, valByte, keySize, entry.valueSize, uint32(entry.timestamp), entrySize)
			_, err = mergeFp.WriteAt(buf, mergeOffset)
			mergeOffset += int64(entrySize)
			if err != nil {
				continue
			}

			b.index.put(string(keyByte), entry)
		}

		b.mu.Lock()
		fp.Close()
		mergeFp.Close()
		log.Printf("Replace old file:'%s' to new merge file: '%s'", oldFilePath, mergeFileName)
		os.Rename(mergeFileName, oldFilePath)
		b.mu.Unlock()
	}
}
