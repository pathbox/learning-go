// Buffer is an binary buffer handler used in emit args. All buffers will be
// sent as binary in the transport layer.

type Buffer {
	Data []byte
	num uint
	isBinary bool
}

// 构造并返回json字符串数据存在buf中
func (b *Buffer) encodeText(buf *bytes.Buffer) error {
	buf.WriteString("{\"type\":\"Buffer\",\"data\":[")
	for i, d := range b.Data {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(int(d)))
	}
	buf.WriteString("]}") // 结束
	return nil
}

func (a *Buffer) encodeBinary(buf *bytes.Buffer) error {
	buf.WriteString("{\"_placeholder\":true,\"num\":")
	buf.WriteString(strconv.FormatUint(a.num, 10))
	buf.WriteString("}")
	return nil
}

// UnmarshalJSON unmarshals from JSON.
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

// MarshalJSON marshals to JSON.
func (a Buffer) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if err := a.marshalJSONBuf(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (a *Buffer) marshalJSONBuf(buf *bytes.Buffer) error {
	encode := a.encodeText
	if a.isBinary {
		encode = a.encodeBinary
	}
	return encode(buf)
}

/*
这段代码来自go-socket中，对json数据进行构造，然后存储到作为参数传入的 buf *bytes.Buffer中，
之后从该buf 中再读取数据。

直接使用 buf.WriteString， 写 构造的json字符串数据，减少了字符串和byte之间的转换

其实也可以定义struct 实例， 然后用json进行序列化，再存到buf， 但这样就多出了序列化的操作，也许性能没有使用上面的方式好
*/