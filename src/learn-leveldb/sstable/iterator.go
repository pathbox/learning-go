package sstable

import (
	"github.com/merlin82/leveldb/internal"
	"github.com/merlin82/leveldb/sstable/block"
)

type Iterator struct {
	table *SsTable
	dataBlockHandle BlockHandle
	dataIter        *block.Iterator
	indexIter       *block.Iterator
}