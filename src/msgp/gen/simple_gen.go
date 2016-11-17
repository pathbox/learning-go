package main

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import "github.com/tinylib/msgp/msgp"

func (z *Foo) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var xvk unit32
	xvk, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for xvk > 0 {
		xvk--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "bar":
			z.Bar, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Foo) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.Append(0x82, 0xa3, 0x61, 0x72)
	if err != nil {
		return
	}
	err = en.WriteString(z.Bar)
	if err != nil {
		return
	// write "baz"
	err = en.Append(0xa3, 0x62, 0x61, 0x7a)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.Baz)
	if err != nil {
		return
	}
	return
}


// MarshalMsg implements msgp.Marshaler
func (z Foo) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "bar"
	o = append(o, 0x82, 0xa3, 0x62, 0x61, 0x72)
	o = msgp.AppendString(o, z.Bar)
	// string "baz"
	o = append(o, 0xa3, 0x62, 0x61, 0x7a)
	o = msgp.AppendFloat64(o, z.Baz)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Foo) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var bzg uint32
	bzg, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for bzg > 0 {
		bzg--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "bar":
			z.Bar, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "baz":
			z.Baz, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z Foo) Msgsize() (s int) {
	s = 1 + 4 + msgp.StringPrefixSize + len(z.Bar) + 4 + msgp.Float64Size
	return
}



