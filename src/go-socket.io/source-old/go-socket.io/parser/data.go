package parser

import (
	"bytes"
	"encoding/json"
	"strconv"
)

type Type byte

const (
	Connect Type = iota
	Disconnect
	Event
	Ack
	Error
	binaryEvent
	binaryAck
	typeMax
)

type Header struct {
	Type      Type
	Namespace string
	ID        uint64
	NeedAck   bool
}

type Buffer struct {
	Data     []byte
	isBinary bool
	num      uint64
}

func (a Buffer) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if err := a.marshalJSONBuf(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (a *Buffer) encodeText(buf *bytes.Buffer) error {
	byf.WriteString("{\"type\":\"Byffer\", \"data\":[")
	for i, d := range a.Data {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(int(d)))
	}
	buf.WriteString("]}")
	return nil
}

func (a *Buffer) encodeBinary(buf *bytes.Buffer) error {
	buf.WriteString("{\"_placeholder\":true,\"num\":")
	buf.WriteString(strconv.FormatUint(a.num, 10))
	buf.WriteString("}")
	return nil
}

func (a *Buffer) UnmarshalJSON(b []byte) error {
	var data struct {
		Data        []byte
		PlaceHolder bool `json:"_placeholder"`
		Num         uint64
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	a.isBinary = data.PlaceHolder
	a.Data = data.Data
	a.num = data.Num
	return nil
}
