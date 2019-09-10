package block

import (
	"bytes"
	"encoding/binary"

	"./internal"
)

type BlockBuilder struct {
	buf     bytes.Buffer
	counter uint32
}

func (blockBuilder *BlockBuilder) Reset() {
	blockBuilder.counter = 0
	blockBuilder.buf.Reset()
}

func (blockBuilder *BlockBuilder) Add(item *internal.InternalKey) error {
	blockBuilder.counter++
	return item.EncodeTo(&blockBuilder.buf)
}

func (blockBuilder *BlockBuilder) Finish() []byte {
	binary.Write(&blockBuilder.buf, binary.LittleEndian, blockBuilder.counter)
	return blockBuilder.buf.Bytes()
}

func (blockBuilder *BlockBuilder) CurrentSizeEstimate() int {
	return blockBuilder.buf.Len()
}

func (blockBuilder *BlockBuilder) Empty() bool {
	return blockBuilder.buf.Len() == 0
}
