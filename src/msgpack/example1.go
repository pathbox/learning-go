package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack"
)

type EventTime struct {
	time.Time
}

var _ msgpack.Marshaler = (*EventTime)(nil)
var _ msgpack.Unmarshaler = (*EventTime)(nil)

func init() {
	msgpack.RegisterExt(0, (*EventTime)(nil))
}

func (tm *EventTime) MarshalMsgpack() ([]byte, error) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b, uint32(tm.Unix())) // 将 tm.Unix()序列化为byte 存到 b，刚好8byte长度

	binary.BigEndian.PutUint32(b[4:], uint32(tm.Nanosecond()))
	return b, nil

}

func (tm *EventTime) UnmarshalMsgpack(b []byte) error {
	if len(b) != 8 {
		return fmt.Errorf("invalid data length: got %d, wanted 8", len(b))
	}

	sec := binary.BigEndian.Uint32(b) // 将byte 转为Uint32 类型值
	usec := binary.BigEndian.Uint32(b[4:])
	tm.Time = time.Unix(int64(sec), int64(usec))
	return nil

}

func main() {
	b, err := msgpack.Marshal(&EventTime{time.Unix(1534567890, 123)})

	if err != nil {
		panic(err)
	}

	var v interface{}
	err = msgpack.Unmarshal(b, &v)
	if err != nil {
		panic(err)
	}

	fmt.Println(v.(*EventTime).UTC())

	tm := new(EventTime)
	err = msgpack.Unmarshal(b, tm)
	if err != nil {
		panic(err)
	}

	fmt.Println(tm.UTC())
}
