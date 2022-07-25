package bitcask

import (
	"encoding/binary"
	"hash/crc32"
)

type entry struct {
	fileID      uint32
	valueSize   uint32
	valueOffset uint64
	timestamp   uint64
}

const HeaderSize = 16

func newEntry(fid, valueSize uint32, valueOffset, timestamp uint64) *entry {
	return &entry{
		fileID:      fid,
		valueSize:   valueSize,
		valueOffset: valueOffset,
		timestamp:   timestamp,
	}
}

func encode(key, value []byte, keySize, valueSize, ts, entrySize uint32) ([]byte, error) {
	// crc32 | timestamp | keySize | valueSize | key | value
	// 4	 | 4		 | 4	   | 4         | 4   | 4
	buf := make([]byte, entrySize)
	binary.BigEndian.PutUint32(buf[4:8], ts)
	binary.BigEndian.PutUint32(buf[8:12], keySize)
	binary.BigEndian.PutUint32(buf[12:16], valueSize)
	copy(buf[HeaderSize:HeaderSize+keySize], key)
	copy(buf[HeaderSize+keySize:HeaderSize+keySize+valueSize], value)

	c32 := crc32.ChecksumIEEE(buf[4:])
	binary.BigEndian.PutUint32(buf[0:4], c32)

	return buf, nil
}

func getSize(keySize, valueSize uint32) uint32 {
	return HeaderSize + keySize + valueSize
}
