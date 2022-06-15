// Open Source: MIT License
// Author: Leon Ding <ding@ibyte.me>
// Date: 2022/2/26 - 10:32 PM - UTC/GMT+08:00

package bottle

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Bottle storage engine components
var (

	// Root Data storage directory
	Root = ""

	// Currently writable file
	active *os.File

	// Concurrent lock
	mutex sync.RWMutex

	// Global indexes
	index map[uint64]*record // 索引其实直接就是一个golang的 map

	// Current data file version
	dataFileVersion int64 = 0

	// Old data file mapping
	fileList map[int64]*os.File

	// Data file name extension
	dataFileSuffix = ".data"

	// index file name extension
	indexFileSuffix = ".index"

	// FRW Read-only Opens a file in write - only mode
	FRW = os.O_RDWR | os.O_APPEND | os.O_CREATE

	// FR Open the file in read-only mode
	FR = os.O_RDONLY

	// Perm Default file operation permission
	Perm = os.FileMode(0750)

	// Default max file size
	// 2 << 8 = 512 << 20 = 536870912 kb
	defaultMaxFileSize int64 = 2 << 8 << 20

	// Data recovery triggers the merge threshold
	totalDataSize int64 = 2 << 8 << 20 << 1 // 1GB

	// Default configuration file format
	defaultConfigFileSuffix = ".yaml"

	// HashedFunc Default Hashed function
	HashedFunc Hashed

	// Secret encryption key
	Secret = []byte("ME:QQ:2420498526")

	// itemPadding binary encoding header padding
	itemPadding uint32 = 20

	// Global data encoder
	encoder *Encoder

	// Write file offset
	writeOffset uint32 = 0

	// Index folder
	indexDirectory string

	// Data folder
	dataDirectory string
)

// Higher-order function blocks
var (

	// Opens a file by specifying a mode
	openDataFile = func(flag int, dataFileIdentifier int64) (*os.File, error) {
		return os.OpenFile(dataSuffixFunc(dataFileIdentifier), flag, Perm)
	}

	// Builds the specified file name extension
	dataSuffixFunc = func(dataFileIdentifier int64) string {
		return fmt.Sprintf("%s%d%s", dataDirectory, dataFileIdentifier, dataFileSuffix)
	}

	// Opens a file by specifying a mode
	openIndexFile = func(flag int, dataFileIdentifier int64) (*os.File, error) {
		return os.OpenFile(indexSuffixFunc(dataFileIdentifier), flag, Perm)
	}

	// Builds the specified file name extension
	indexSuffixFunc = func(dataFileIdentifier int64) string {
		return fmt.Sprintf("%s%d%s", indexDirectory, dataFileIdentifier, indexFileSuffix)
	}
)

// record Mapping Data Record
type record struct {
	FID        int64  // data file id
	Size       uint32 // data record size
	Offset     uint32 // data record offset
	Timestamp  uint32 // data record create timestamp
	ExpireTime uint32 // data record expire time
}

func Open(opt Option) error {
	opt.Validation()

	initialize()

	if ok, err := pathExists(Root); ok {
		// The directory has recovered data. Procedure
		return recoverData()
	} else if err != nil {
		// If there is an error, the file is not a directory or is invalid
		panic("The current path is invalid!!!")
	}
	// Create folder if it does not exist
	if err := os.MkdirAll(dataDirectory, Perm); err != nil {
		panic("Failed to create a working directory!!!")
	}

	if err := os.MkdirAll(indexDirectory, Perm); err != nil {
		panic("Failed to create a working directory!!!")
	}

	// Once the directory is created, you can create active files to write data
	return createActiveFile()
}

// Load through a configuration file
func Load(file string) error {
	if path.Ext(file) != defaultConfigFileSuffix {
		panic("the current configuration file format is invalid")
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(file)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	var opt Option
	if err := v.Unmarshal(&opt); err != nil {
		return err
	}

	return Open(opt)
}

// Action Operation add-on
type Action struct {
	TTL time.Time // Survival time
}

// TTL You can set a timeout for the key in seconds
func TTL(second uint32) func(action *Action) {
	return func(action *Action) {
		action.TTL = time.Now().Add(time.Duration(second) * time.Second)
	}
}

// Put Add key-value data to the storage engine
// actionFunc You can set the timeout period
func Put(key, value []byte, actionFunc ...func(action *Action)) (err error) {
	var (
		action Action
		size   int
	)

	if len(actionFunc) > 0 {
		for _, fn := range actionFunc {
			fn(&action)
		}
	}

	fileInfo, _ := active.Stat()

	if fileInfo.Size() >= defaultMaxFileSize {
		if err := closeActiveFile(); err != nil {
			return err
		}

		if err := createActiveFile(); err != nil {
			return err
		}
	}

	sum64 := HashedFunc.Sum64(key)

	mutex.Lock()
	defer mutex.Unlock()

	timestamp := time.Now().Unix()

	if size, err = encoder.Write(NewItem(key, value, uint64(timestamp)), active); err != nil {
		return err
	}

	index[sum64] = &record{
		FID:        dataFileVersion,
		Size:       uint32(size),
		Offset:     writeOffset,
		Timestamp:  uint32(timestamp),
		ExpireTime: uint32(action.TTL.Unix()),
	}

	writeOffset += uint32(size)

	return nil
}

// Get gets the data object for the specified key
func Get(key []byte) (data *Data) {
	data = &Data{}
	mutex.RLock()
	defer mutex.RUnlock()

	sum64 := HashedFunc.Sum64(key)

	if index[sum64] == nil {
		data.Err = errors.New("the current key does not exist")
		return
	}

	if index[sum64].ExpireTime <= uint32(time.Now().Unix()) {
		data.Err = errors.New("the current key has expired")
		return
	}

	item, err := encoder.Read(index[sum64])
	if err != nil {
		data.Err = err
		return
	}
	data.Item = item
	return
}

// Remove removes specified data from storage
func Remove(key []byte) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(index, HashedFunc.Sum64(key))
}

// Close shut down the storage engine and flush the data
func Close() error {
	mutex.Lock()
	defer mutex.Unlock()

	if err := active.Sync(); err != nil {
		return err
	}

	for _, file := range fileList {
		if err := file.Close(); err != nil {
			return err
		}
	}

	return saveIndexToFile()
}

// Create a new active file
func createActiveFile() error {
	mutex.Lock()
	defer mutex.Unlock()

	// Initialize writable file offsets and file identifiers
	writeOffset = 0
	dataFileVersion++

	if file, err := openDataFile(FRW, dataFileVersion); err == nil {
		active = file
		fileList[dataFileVersion] = active
		return nil
	}

	return errors.New("failed to create writable data file")
}

func closeActiveFile() error {
	mutex.Lock()
	defer mutex.Unlock()

	if err := active.Sync(); err != nil {
		return err
	}

	if err := active.Close(); err != nil {
		return err
	}

	// Set the previous writable file to read-only
	if file, err := openDataFile(FR, dataFileVersion); err == nil {
		fileList[dataFileVersion] = file
		return nil
	}

	return errors.New("error opening write only file")
}

// Initialize storage engine components
func initialize() {
	if HashedFunc == nil {
		HashedFunc = DefaultHashFunc()
	}
	if encoder == nil {
		encoder = DefaultEncoder()
	}
	if index == nil {
		index = make(map[uint64]*record)
	}

	// By default, five file descriptors are mounted
	fileList = make(map[int64]*os.File, 5)
}

// Memory index file item encoding used
// The size of 288 - bit
type indexItem struct {
	idx uint64
	*record
}

// Save index files to the data directory
func saveIndexToFile() (err error) {
	var file *os.File
	defer func() {
		if err := file.Sync(); err != nil {
			return
		}
		if err := file.Close(); err != nil {
			return
		}
	}()

	var channel = make(chan indexItem, 1024)

	go func() {
		for sum64, record := range index {
			channel <- indexItem{
				idx:    sum64,
				record: record,
			}
		}
		close(channel)
	}()

	if file, err = openIndexFile(FRW, time.Now().Unix()); err != nil {
		return
	}

	for v := range channel {
		if _, err = encoder.WriteIndex(v, file); err != nil {
			return
		}
	}

	return
}

func recoverData() error {

	if dataTotalSize() >= totalDataSize {
		// Trigger merger
		if err := migrate(); err != nil {
			return err
		}
	}

	// Find the last data file and see if it's full
	if file, err := findLatestDataFile(); err == nil {
		info, _ := file.Stat()
		if info.Size() >= defaultMaxFileSize {
			if err := createActiveFile(); err != nil {
				return err
			}
			// When the data is full, a new writable file is created and an index is built
			return buildIndex()
		}
		// If the data file was not full last time
		// it is set to writable and the writable offset is calculated
		active = file
		if offset, err := file.Seek(0, os.SEEK_END); err == nil {
			writeOffset = uint32(offset)
		}
		return buildIndex()
	}

	return errors.New("failed to restore data")
}

// Trigger data file merge Dirty data merge
func migrate() error {
	// Load indexes and load data
	if err := buildIndex(); err != nil {
		return err
	}

	// Get the latest version of the data
	version()

	var (
		offset       uint32
		file         *os.File
		fileInfo     os.FileInfo
		excludeFiles []int64
		activeItem   = make(map[uint64]*Item, len(index))
	)

	dataFileVersion++
	// Create the target data file for migration
	file, _ = openDataFile(FRW, dataFileVersion)
	excludeFiles = append(excludeFiles, dataFileVersion)
	// Get the migration file status
	fileInfo, _ = file.Stat()

	// Migrate active recordable
	for idx, rec := range index {
		item, err := encoder.Read(rec)

		if err != nil {
			return err
		}

		activeItem[idx] = item
	}

	for idx, item := range activeItem {
		// Check whether the migration file threshold is reached at each turn
		if fileInfo.Size() >= defaultMaxFileSize {
			// Close and set too read-only to put into map
			if err := file.Sync(); err != nil {
				return err
			}
			if err := file.Close(); err != nil {
				return err
			}

			// The update operation
			dataFileVersion++
			excludeFiles = append(excludeFiles, dataFileVersion)

			file, _ = openDataFile(FRW, dataFileVersion)
			fileInfo, _ = file.Stat()
			offset = 0
		}

		// Write the original content to the new file
		size, err := encoder.Write(item, file)

		if err != nil {
			return err
		}

		// Update the new file ID and offset
		index[idx].FID = dataFileVersion
		index[idx].Size = uint32(size)
		index[idx].Offset = offset

		offset += uint32(size)
	}

	// Clear deleted data
	fileInfos, err := ioutil.ReadDir(dataDirectory)

	if err != nil {
		return err
	}

	// Filter out data files that have been migrated
	for _, info := range fileInfos {
		fileName := fmt.Sprintf("%s%s", dataDirectory, info.Name())
		for _, excludeFile := range excludeFiles {
			if fileName != dataSuffixFunc(excludeFile) {
				if err := os.Remove(fileName); err != nil {
					return err
				}
			}
		}
	}

	// After the migration, save the latest index file
	return saveIndexToFile()
}

func buildIndex() error {

	if err := readIndexItem(); err != nil {
		return err
	}

	// Find the data file from the index and open the file descriptor
	for _, record := range index {
		// https://stackoverflow.com/questions/37804804/too-many-open-file-error-in-golang
		if fileList[record.FID] == nil {
			file, err := openDataFile(FR, record.FID)
			if err != nil {
				return err
			}
			// Open the original data file
			fileList[record.FID] = file
		}
	}

	return nil
}

// Find the latest data files in the index folder
func findLatestIndexFile() (*os.File, error) {
	files, err := ioutil.ReadDir(indexDirectory)

	if err != nil {
		return nil, err
	}

	var indexes []fs.FileInfo

	for _, file := range files {
		if path.Ext(file.Name()) == indexFileSuffix {
			indexes = append(indexes, file)
		}
	}

	var ids []int

	for _, info := range indexes {
		id := strings.Split(info.Name(), ".")[0]
		i, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, i)
	}

	sort.Ints(ids)

	return openIndexFile(FR, int64(ids[len(ids)-1]))
}

// Read index file contents into memory index
func readIndexItem() error {
	if file, err := findLatestIndexFile(); err == nil {
		defer func() {
			if err := file.Sync(); err != nil {
				return
			}
			if err := file.Close(); err != nil {
				return
			}
		}()

		buf := make([]byte, 36)

		for {
			_, err := file.Read(buf)

			if err != nil && err != io.EOF {
				return err
			}

			if err == io.EOF {
				break
			}

			if err = encoder.ReadIndex(buf); err != nil {
				return err
			}
		}

		return nil
	}

	return errors.New("index reading failed")
}

// Find the latest data file from the data file
func findLatestDataFile() (*os.File, error) {
	version()
	return openDataFile(FRW, dataFileVersion)
}

// Load the data file version number
func version() {
	files, _ := ioutil.ReadDir(dataDirectory)

	var datafiles []fs.FileInfo

	for _, file := range files {
		if path.Ext(file.Name()) == dataFileSuffix {
			datafiles = append(datafiles, file)
		}
	}

	var ids []int

	for _, info := range datafiles {
		id := strings.Split(info.Name(), ".")[0]
		i, _ := strconv.Atoi(id)
		ids = append(ids, i)
	}

	sort.Ints(ids)

	// Reset file counters and writable files and offsets
	dataFileVersion = int64(ids[len(ids)-1])
}

// Calculate all data file sizes from the data folder
func dataTotalSize() int64 {
	files, _ := ioutil.ReadDir(dataDirectory)

	var datafiles []fs.FileInfo

	for _, file := range files {
		if path.Ext(file.Name()) == dataFileSuffix {
			datafiles = append(datafiles, file)
		}
	}

	var totalSize int64

	for _, datafile := range datafiles {
		totalSize += datafile.Size()
	}

	return totalSize
}
