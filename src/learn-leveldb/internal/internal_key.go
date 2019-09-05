package internal

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
)

type ValueType int8

const(
	TypeDeletion ValueType = 0
	TypeValue ValueType = 1
)

type InternalKey struct {
	Seq uint64
	Type ValueType
	UserKey []byte
	UserValue []byte
}

func NewInternalKey(seq uint64, valueType ValueType, key, value []byte) *InternalKey {
	var internalKey InternalKey
	internalKey.Seq = seq
	internalKey.Type = valueType
	internalKey.UserKey = make([]byte, len(key))
	copy(internalKey.UserKey, key)
	internalKey.UserValue = make([]byte, len(value))
	copy(internalKey.UserValue, value) // slice 数据要使用copy深度copy

	return &internalKey
}
func (key *InternalKey) EncodeTo(w io.Writer) error {
	binary.Write(w, binary.LittleEndian, key.Seq)  // LittleEndian 的方式编码
}