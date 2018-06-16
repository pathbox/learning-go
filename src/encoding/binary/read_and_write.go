import (
	"errors"
	"io"
	"reflect"
)

func Read(r io.Reader, order ByteOrder, data interface{}) error {
	// Fast path for basic types and slices.
	if n := intDataSize(data); n != 0 {
		var b [8]byte
		var bs []byte
		if n > len(b) {
			bs = make([]byte, n)
		} else {
			bs = b[:n]
		}
		if _, err := io.ReadFull(r, bs); err != nil {
			return err
		}
		switch data := data.(type) {
		case *bool:
			*data = b[0] != 0
		case *int8:
			*data = int8(b[0])
		case *uint8:
			*data = b[0]
		case *int16:
			*data = int16(order.Uint16(bs))
		case *uint16:
			*data = order.Uint16(bs)
		case *int32:
			*data = int32(order.Uint32(bs))
		case *uint32:
			*data = order.Uint32(bs)
		case *int64:
			*data = int64(order.Uint64(bs))
		case *uint64:
			*data = order.Uint64(bs)
		case []bool:
			for i, x := range bs { // Easier to loop over the input for 8-bit values.
				data[i] = x != 0
			}
		case []int8:
			for i, x := range bs {
				data[i] = int8(x)
			}
		case []uint8:
			copy(data, bs)
		case []int16:
			for i := range data {
				data[i] = int16(order.Uint16(bs[2*i:]))
			}
		case []uint16:
			for i := range data {
				data[i] = order.Uint16(bs[2*i:])
			}
		case []int32:
			for i := range data {
				data[i] = int32(order.Uint32(bs[4*i:]))
			}
		case []uint32:
			for i := range data {
				data[i] = order.Uint32(bs[4*i:])
			}
		case []int64:
			for i := range data {
				data[i] = int64(order.Uint64(bs[8*i:]))
			}
		case []uint64:
			for i := range data {
				data[i] = order.Uint64(bs[8*i:])
			}
		}
		return nil
	}

	// Fallback to reflect-based decoding.
	v := reflect.ValueOf(data)
	size := -1
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
		size = dataSize(v)
	case reflect.Slice:
		size = dataSize(v)
	}
	if size < 0 {
		return errors.New("binary.Read: invalid type " + reflect.TypeOf(data).String())
	}
	d := &decoder{order: order, buf: make([]byte, size)}
	if _, err := io.ReadFull(r, d.buf); err != nil {
		return err
	}
	d.value(v)
	return nil
}

// Write writes the binary representation of data into w.
// Data must be a fixed-size value or a slice of fixed-size
// values, or a pointer to such data.
// Boolean values encode as one byte: 1 for true, and 0 for false.
// Bytes written to w are encoded using the specified byte order
// and read from successive fields of the data.
// When writing structs, zero values are written for fields
// with blank (_) field names.
func Write(w io.Writer, order ByteOrder, data interface{}) error {
	// Fast path for basic types and slices.
	if n := intDataSize(data); n != 0 {
		var b [8]byte
		var bs []byte
		if n > len(b) {
			bs = make([]byte, n)
		} else {
			bs = b[:n]
		}
		switch v := data.(type) {
		case *bool:
			if *v {
				b[0] = 1
			} else {
				b[0] = 0
			}
		case bool:
			if v {
				b[0] = 1
			} else {
				b[0] = 0
			}
		case []bool:
			for i, x := range v {
				if x {
					bs[i] = 1
				} else {
					bs[i] = 0
				}
			}
		case *int8:
			b[0] = byte(*v)
		case int8:
			b[0] = byte(v)
		case []int8:
			for i, x := range v {
				bs[i] = byte(x)
			}
		case *uint8:
			b[0] = *v
		case uint8:
			b[0] = v
		case []uint8:
			bs = v
		case *int16:
			order.PutUint16(bs, uint16(*v))
		case int16:
			order.PutUint16(bs, uint16(v))
		case []int16:
			for i, x := range v {
				order.PutUint16(bs[2*i:], uint16(x))
			}
		case *uint16:
			order.PutUint16(bs, *v)
		case uint16:
			order.PutUint16(bs, v)
		case []uint16:
			for i, x := range v {
				order.PutUint16(bs[2*i:], x)
			}
		case *int32:
			order.PutUint32(bs, uint32(*v))
		case int32:
			order.PutUint32(bs, uint32(v))
		case []int32:
			for i, x := range v {
				order.PutUint32(bs[4*i:], uint32(x))
			}
		case *uint32:
			order.PutUint32(bs, *v)
		case uint32:
			order.PutUint32(bs, v)
		case []uint32:
			for i, x := range v {
				order.PutUint32(bs[4*i:], x)
			}
		case *int64:
			order.PutUint64(bs, uint64(*v))
		case int64:
			order.PutUint64(bs, uint64(v))
		case []int64:
			for i, x := range v {
				order.PutUint64(bs[8*i:], uint64(x))
			}
		case *uint64:
			order.PutUint64(bs, *v)
		case uint64:
			order.PutUint64(bs, v)
		case []uint64:
			for i, x := range v {
				order.PutUint64(bs[8*i:], x)
			}
		}
		_, err := w.Write(bs)
		return err
	}

	// Fallback to reflect-based encoding.
	v := reflect.Indirect(reflect.ValueOf(data))
	size := dataSize(v)
	if size < 0 {
		return errors.New("binary.Write: invalid type " + reflect.TypeOf(data).String())
	}
	buf := make([]byte, size)
	e := &encoder{order: order, buf: buf}
	e.value(v)
	_, err := w.Write(buf)
	return err
}
