package block

import (
	"bytes"
	"encoding/binary"

	"./internal"
)

type Block struct {
	items []internal.InternalKey
}

func New(p []byte) *Block {
	var block Block
	data := bytes.NewBuffer(p)
	counter := binary.LittleEndian.Uint32(p[len(p)-4:])

	for i := uint32(0); i < counter; i++{
		var item internal.InternalKey
		err := item.DecodeFrom(data)
		if err != nil {
			return nil
		}
		block.items = append(block.items, item)
	}

	return &block
}

func (block *Block) NewIterator() *Iterator {
	return &Iterator{block: block}
}